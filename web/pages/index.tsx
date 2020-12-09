import React, { ReactElement } from "react";
import Libp2p from "libp2p";
import Head from "next/head";
import { useState } from "react";
import { startNode } from "../services/libp2p";
import { createOrGetPeerID } from "../services/peerid";

export default function Home(): ReactElement {
  const [libp2p, setLibp2p] = useState<Libp2p | undefined>(undefined);
  const [ID, setID] = useState<string>("Loading...");

  const connect = async () => {
    const peerID = await createOrGetPeerID();
    console.log({ peerID });
    setID(peerID.toB58String());
    const _libp2p = await startNode(peerID);
    _libp2p.connectionManager.on("peer:connect", (connection) =>
      console.log("Connected peer: ", connection)
    );
    _libp2p.connectionManager.on("peer:disconnect", (connection) =>
      console.log("Disconnected peer: ", connection)
    );
    _libp2p.on("peer:discovery", (peerId) =>
      console.log("Peer discovery: ", peerId)
    );

    _libp2p.pubsub.subscribe("chuj123");

    _libp2p.pubsub.on("chuj123", (msg) => console.log({ msg }));

    setLibp2p(_libp2p);
  };

  if (libp2p === undefined) {
    return (
      <>
        Not connected
        <br />
        <button
          onClick={async () => {
            try {
              await connect();
            } catch (e) {
              console.log("Error when connecting: ", e);
              alert(e.message);
            }
          }}
        >
          Connect
        </button>
      </>
    );
  }

  const publishMsg = async () => {
    const textEncoder = new TextEncoder();
    const data: Uint8Array = textEncoder.encode(
      `c: ${Math.round(Math.random())}`
    );
    const res = libp2p.pubsub.publish("chuj123", data);
  };

  console.log(libp2p);
  return (
    <div>
      <Head>
        <title>Create Next App</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      Connected
      <br />
      <button onClick={() => console.log(libp2p)}>Log list of peers</button>
      <br />
      <button onClick={publishMsg}>Publish message</button>
      <br />
      <br />
      ID: {ID}
      <br />
    </div>
  );
}
