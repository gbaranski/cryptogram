package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4001")
	if err != nil {
		log.Panicln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting LibP2P host")
	host, err := libp2p.New(ctx, libp2p.ListenAddrs(addr))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("LibP2P host started")
	log.Println("Addresses: ")
	for i, addr := range host.Addrs() {
		log.Println(i, " - ", addr)
	}
	_, err = dht.New(ctx, host)
	select {}
}
