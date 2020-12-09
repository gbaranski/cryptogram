import Libp2p, { Options } from "libp2p";
import Websockets from "libp2p-websockets";
import WebrtcStar from "libp2p-webrtc-star";
import Mplex from "libp2p-mplex";
import { NOISE } from "libp2p-noise";
import Bootstrap from "libp2p-bootstrap";
import KadDHT from "libp2p-kad-dht";
import Gossipsub from "libp2p-gossipsub";

import PeerId from "peer-id";

export const startNode = async (peerID: PeerId): Promise<Libp2p> => {
  const options: Options = {
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
        enabled: true,
        randomWalk: {
          enabled: true,
        },
      },
    },
  };

  const libp2p: Libp2p = await Libp2p.create(options);
  await libp2p.start();
  return libp2p;
};
