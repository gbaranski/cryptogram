package cli

import (
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

func executeCommand(cmd string, api *node.API) {
	args := strings.Split(cmd, " ")

	switch args[0] {
	case "":
		return
	case "exit", "quit":
		fmt.Println("Goodbye")
		os.Exit(0)
	case "id":
		fmt.Println((*api.Host).ID())
	case "peers":
		for i, p := range (*api.Host).Network().Peers() {
			fmt.Println(i, " - ", p)
		}
	default:
		fmt.Printf("Unhandled command \"%s\"\n", args[0])
	}

}

// Executor used to execute commands in CLI
func Executor(cmd string, api *node.API) {
	if strings.HasPrefix(cmd, "/") {
		executeCommand(strings.TrimPrefix(cmd, "/"), api)
		return
	}

	fmt.Println("Sending message: ", cmd)
}
