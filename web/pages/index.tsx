import Head from "next/head";
import { useState } from "react";
import { startNode } from "../services/libp2p";
import { createOrGetPeerID } from "../services/peerid";

export default function Home() {
  const [libp2p, setLibp2p] = useState<LibP2p | undefined>(undefined);
  const [ID, setID] = useState<string>("Loading...");

  const connect = async () => {
    const peerID = await createOrGetPeerID();
    const _libp2p = await startNode(peerID);
    setLibp2p(_libp2p);
  };

  if (libp2p === undefined) {
    return (
      <>
        Not connected
        <br />
        <button onClick={connect}>Connect</button>
      </>
    );
  }

  console.log(libp2p);
  return (
    <div>
      <Head>
        <title>Create Next App</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      Connected
      <br />
      {/* <button onClick={() => console.log(IPFS.pubsub.peers("chuj123"))}>
        Log list of peers
      </button>
      <br />
      <button onClick={() => IPFS.pubsub.publish("chuj123", "hellothere")}>
        Send message
      </button>
      <br />
      <button
        onClick={() =>
          IPFS.pubsub.subscribe("chuj123", (msg) => console.log({ msg }))
        }
      >
        Subscribe
      </button> */}
      <br />
      ID: {ID}
      <br />
    </div>
  );
}
