package chat

import (
	"context"

	"github.com/gbaranski/cryptogram/cli/misc"
	"github.com/gbaranski/cryptogram/cli/node"
)

// Chat hold current chat state
type Chat struct {
	Context   *context.Context
	API       *node.API
	MsgSender *MessageSender
}

// CreateChat creates chat
func CreateChat(context context.Context, config *misc.Config, API *node.API) *Chat {
	chat := &Chat{
		Context: &context,
		API:     API,
		MsgSender: &MessageSender{
			PeerID:   (*API.Host).ID(),
			Nickname: *config.Nickname,
		},
	}
	return chat
}
