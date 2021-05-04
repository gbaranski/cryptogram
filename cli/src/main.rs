use cryptogram_discovery::mdns;
use std::time::Duration;
use structopt::StructOpt;
use tokio::time;

#[derive(StructOpt, Debug)]
#[structopt(name = "Cryptogram")]
struct Opt {
    // The number of occurrences of the `v/verbose` flag
    /// Verbose mode (-v, -vv, -vvv, etc.)
    #[structopt(short, long, parse(from_occurrences))]
    verbose: u8,
}

async fn mdns_announce_loop(discovery: mdns::Discovery) {
    let mut interval = time::interval(Duration::from_secs(10));
    loop {
        interval.tick().await;
        discovery.announce().await.expect("failed announcing");
        log::debug!("Announced service via mDNS");
    }
}

async fn mdns_poll_loop(discovery: mdns::Discovery) {
    loop {
        let socket_addr = discovery.poll().await.unwrap();
        println!("Discovered new client with socket addr: {:?}", socket_addr);
    }
}

#[tokio::main]
async fn main() {
    let opt = Opt::from_args();
    loggerv::init_with_verbosity(opt.verbose as u64).expect("failed initializing verbosity");
    println!("Hello, world!");
    let mdns_discovery = mdns::Discovery::new().await.unwrap();
    tokio::select! {
        _ = mdns_announce_loop(mdns_discovery.clone()) => {
        },
        _ = mdns_poll_loop(mdns_discovery.clone()) => {
        }
    };
}
