package chat

import "github.com/libp2p/go-libp2p-core/peer"

// MessageSender sender field of {Message}
type MessageSender struct {
	PeerID   peer.ID `json:"peerID"`
	Nickname string  `json:"nickname"`
}

// Message is message
type Message struct {
	Text   string        `json:"text"`
	Sender MessageSender `json:"sender"`
}
