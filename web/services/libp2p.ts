import Libp2p from "libp2p";
import Websockets from "libp2p-websockets";
import WebrtcStar from "libp2p-webrtc-star";
import Mplex from "libp2p-mplex";
import { NOISE } from "libp2p-noise";
import Bootstrap from "libp2p-bootstrap";
import KadDHT from "libp2p-kad-dht";
import Gossipsub from "libp2p-gossipsub";

// import "@polkadot/ts/libp2p";
// import "@polkadot/ts/libp2p-bootstrap";
// import "@polkadot/ts/libp2p-crypto";
// import "@polkadot/ts/libp2p-kad-dht";
// import "@polkadot/ts/libp2p-mdns";
// import "@polkadot/ts/libp2p-mplex";
// import "@polkadot/ts/libp2p-secio";
// import "@polkadot/ts/libp2p-spdy";
// import "@polkadot/ts/libp2p-tcp";
// import "@polkadot/ts/libp2p-webrtc-direct";
// import "@polkadot/ts/libp2p-webrtc-star";
// import "@polkadot/ts/libp2p-websocket-star";
// import "@polkadot/ts/libp2p-websockets";
import PeerId from "peer-id";

export const startNode = async (peerID: PeerId): Promise<Libp2p> => {
  const options: LibP2p.Options = {
    peerId: peerID,
    addresses: {
      listen: [
        // Add the signaling server multiaddr
        "/ip4/127.0.0.1/tcp/15555/ws/p2p-webrtc-star",
      ],
    },
    modules: {
      transport: [Websockets, WebrtcStar],
      streamMuxer: [Mplex],
      // @ts-ignore
      connEncryption: [NOISE],
      peerDiscovery: [Bootstrap],
      dht: KadDHT,
      pubsub: Gossipsub,
    },
    config: {
      peerDiscovery: {
        bootstrap: {
          list: [
            "/ip4/127.0.0.1/tcp/63786/ws/p2p/QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d",
          ],
        },
      },
      dht: {
        // @ts-ignore
        enabled: true,
        randomWalk: {
          enabled: true,
        },
      },
    },
  };

  // @ts-ignore
  const libp2p: Libp2p = await Libp2p.create(options);
  await libp2p.start();
  return libp2p;
};
