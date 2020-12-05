package repl

import (
	"context"
	"fmt"
	"os"
	"strings"

	iface "github.com/ipfs/interface-go-ipfs-core"
)

// Executor used to execute commands in CLI
func Executor(s string, ipfs iface.CoreAPI) {
	s = strings.TrimSpace(s)
	switch s {
	case "":
		return
	case "exit", "quit":
		fmt.Println("Goodbye")
		os.Exit(0)
	case "peers":
		fmt.Print("Peers: ")
		fmt.Println(ipfs.PubSub().Peers(context.TODO()))
	default:
		fmt.Printf("Unhandled command \"%s\"\n", s)
	}
}
