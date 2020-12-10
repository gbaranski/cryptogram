package misc

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p-core/network"
)

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			log.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		log.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			log.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			log.Println("Error flushing buffer")
			panic(err)
		}
	}
}

// HandleNetworkStream handles network stream
func HandleNetworkStream(stream network.Stream) {
	log.Println("Got a new stream!")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go readData(rw)
	go writeData(rw)
}
