package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	fmt.Print("\n\nWELCOME TO P2P CHAT APPLICATION\n\n")
	
	portF := flag.Int("p", 0, "port number")
	targetF := flag.String("t", "", "target address")

	flag.Parse()

	if *portF == 0 {
		log.Println("Error: No port number passed.")
		return
	}
	
	node, err := CreateNode(*portF)
	if err != nil {
		log.Panic(err)
	}
	defer node.Close()

	if *targetF == "" {
		// if there's no target, create a new stream since this node is the target
		CreateNewStreamWithNode(node, ChatProtocolId)

		fmt.Printf("\n\nRun this in new terminal:\n\n./basic-chat -p %v -t /ip4/127.0.0.1/tcp/%v/p2p/%s\n", *portF + 1, *portF, node.ID())

	} else {
		// else connect to the target node and attach the node to the stream
		readWriter, err := ConnectNodeToStream(node, *targetF, ChatProtocolId)
		if err != nil {
			log.Panic(err)
		}

		go WriteStream(readWriter)
		go ReadStream(readWriter)
	}

	// keep the terminal open
	select {}
}
