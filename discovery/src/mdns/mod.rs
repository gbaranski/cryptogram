use bytes::{BufMut, BytesMut};
use std::net::{IpAddr, Ipv4Addr, Ipv6Addr, SocketAddr, SocketAddrV4};
use std::sync::Arc;
use thiserror::Error;
use tokio::net::UdpSocket;

mod dns;

pub const SERVICE_NAME: &'static str = "_cryptogram._tcp.local";
pub const MULTICAST_PORT: u16 = 5353;
pub const MULTICAST_IPV4: Ipv4Addr = Ipv4Addr::new(224, 0, 0, 123);
pub const MULTICAST_IPV6: Ipv6Addr = Ipv6Addr::new(0xFF02, 0, 0, 0, 0, 0, 0, 0x0123);

#[derive(Error, Debug)]
pub enum DiscoveryError {
    #[error("IO Error occured: `{source}`")]
    IOError {
        #[from]
        source: std::io::Error,
    },
}

#[derive(Clone)]
pub struct Discovery {
    socket_rx: Arc<UdpSocket>,
    socket_tx: Arc<UdpSocket>,
}

impl Discovery {
    fn new_socket_rx() -> Result<UdpSocket, DiscoveryError> {
        use socket2::{Domain, Protocol, Socket, Type};
        let addr = SocketAddr::new(Ipv4Addr::UNSPECIFIED.into(), MULTICAST_PORT).into();
        let socket = Socket::new(Domain::IPV4, Type::DGRAM, Some(Protocol::UDP))?;
        socket.set_reuse_address(true)?;
        socket.bind(&addr)?;
        socket.set_multicast_loop_v4(true)?;
        socket.set_multicast_ttl_v4(255)?;
        let std_socket = socket.into();
        let async_socket = UdpSocket::from_std(std_socket)?;
        Ok(async_socket)
    }

    async fn new_socket_tx() -> Result<UdpSocket, DiscoveryError> {
        let addr = SocketAddr::new(Ipv4Addr::UNSPECIFIED.into(), 0);
        let socket = UdpSocket::bind(addr).await?;
        Ok(socket)
    }

    pub async fn new() -> Result<Self, DiscoveryError> {
        let address = SocketAddr::V4(SocketAddrV4::new(Ipv4Addr::UNSPECIFIED, MULTICAST_PORT));
        let socket_rx = Self::new_socket_rx()?;
        let socket_tx = Self::new_socket_tx().await?;

        Ok(Discovery {
            socket_rx: Arc::new(socket_rx),
            socket_tx: Arc::new(socket_tx),
        })
    }

    pub async fn announce(&self) -> Result<(), DiscoveryError> {
        let mut buf = BytesMut::with_capacity(512);
        let id = 1337_u16; // No idea why this exact number
        let qr = true; // A one bit field that specifies whether this message is a query, or a response.

        // Source: https://www2.cs.duke.edu/courses/fall16/compsci356/DNS/DNS-primer.pdf

        buf.put_u16(id);
        buf.put_u16(0x8400); // Answer flag
        buf.put_u16(0x0); // Number of questions
        buf.put_u16(0x1); // Number of answers
        buf.put_u16(0x0); // Number of authorities
        buf.put_u16(0x0); // Number of additionals

        buf.put_u8(b'\0');
        unimplemented!();
    }

    pub async fn poll(&self) -> Result<SocketAddr, DiscoveryError> {
        let mut buf = vec![0; 512];
        let (len, remote_address) = self.socket_rx.recv_from(&mut buf).await?;
        let data = &buf[..len];
        let response = String::from_utf8_lossy(data);
        log::info!("client got data: {}, from: {}", response, remote_address);
        Ok(remote_address)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn test_ipv4_multicast() {
        assert!(MULTICAST_IPV4.is_multicast())
    }
    #[test]
    fn test_ipv6_multicast() {
        assert!(MULTICAST_IPV6.is_multicast())
    }
}
