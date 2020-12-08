import PeerId from "peer-id";

const getPeerID = async (): Promise<PeerId | undefined> => {
  const stringifiedPeerId = localStorage.getItem("peerID");
  if (!stringifiedPeerId) return undefined;
  console.info("Retreived successfully PeerID from LocalStorage");
  const peerIDJson = JSON.parse(stringifiedPeerId);
  return PeerId.createFromJSON(peerIDJson);
};

const updatePeerID = (peerID: PeerId) => {
  const stringifiedPeerId = JSON.stringify(peerID);
  console.info("Updating PeerID in LocalStorage");
  localStorage.setItem("peerID", stringifiedPeerId);
};

export const createOrGetPeerID = async (): Promise<PeerId> => {
  let peerID = await getPeerID();
  if (peerID) return peerID;
  console.info(
    "Couldn't find existing PeerID in LocalStorage, generating new one..."
  );
  peerID = await PeerId.create();
  updatePeerID(peerID);
  return peerID;
};
