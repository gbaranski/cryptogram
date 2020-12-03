package main

import (
	"context"
	"fmt"

	"github.com/c-bata/go-prompt"
	repl "github.com/gbaranski/cryptogram/cli/repl"

	node "github.com/gbaranski/cryptogram/cli/node"
)

func main() {

	p := prompt.New(repl.Executor, repl.Completer, prompt.OptionTitle("cryptogram-cli"), prompt.OptionPrefix(">>> "))
	p.Run()

	return

	fmt.Println("-- Getting an IPFS node running -- ")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println("Spawning node on a temporary repo")
	_, err := node.SpawnEphemeral(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to spawn ephemeral node: %s", err))
	}
	fmt.Println("IPFS node is running")

}
