package repl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gbaranski/cryptogram/cli/node"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

func validateArgs(args []string, max int, min int) bool {
	if len(args) > max {
		fmt.Printf("Unexpected argument %s", args[len(args)-1])
		return false
	} else if len(args) < min {
		fmt.Printf("More arguments required")
		return false
	}
	return true
}

var subscription iface.PubSubSubscription

// Executor used to execute commands in CLI
func Executor(fullCMD string, ipfs node.IPFSNodeAPI) {
	args := strings.Split(fullCMD, " ")

	switch args[0] {
	case "":
		return
	case "exit", "quit":
		fmt.Println("Goodbye")
		os.Exit(0)
	case "id":
		fmt.Printf("Node ID: %s\n", ipfs.Node.Identity.String())
	case "peers":
		peers, err := ipfs.API.PubSub().Peers(context.Background())
		if err != nil {
			fmt.Printf("Error occured when retreiving list of peers %s", err)
		}
		fmt.Println("Peers: ")
		for _, v := range peers {
			fmt.Println(v)

		}
	case "subscriptions":
		fmt.Print("Subscriptions: ")
		fmt.Println(ipfs.API.PubSub().Ls(context.TODO()))
	case "send_message":
		if validateArgs(args, 3, 3) == false {
			return
		}
		node.SendMessage(ipfs, args[1], args[2])
	case "topic":
		switch args[1] {
		case "subscribe":
			if validateArgs(args, 3, 3) == false {
				return
			}
			subscription = node.SubscribeMessages(ipfs, args[2])
		case "unsubscribe":
			fmt.Println("Not implemented yet")
		default:
			fmt.Printf("Unhandled command \"%s\"\n", args[0])
		}
	case "listen":
		go node.ListenToMessages(subscription)
	default:
		fmt.Printf("Unhandled command \"%s\"\n", args[0])
	}
}
