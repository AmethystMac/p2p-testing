package main

import (
	"fmt"
	"log"

	"blockchain-prototype/p2p"
	"blockchain-prototype/stream"
)

func main() {
	
	// _ = flag.Int("p", 0, "port number")
	// targetF := flag.String("t", "", "target address")
	
	node1, err := p2p.CreateNode(5500)
	if err != nil {
		log.Panic(err)
	}
	defer node1.Close()

	p2p.CreateNewStreamWithNode(node1, stream.ProtocolId)

	node2, err := p2p.CreateNode(5600)
	if err != nil {
		log.Panic(err)
	}
	defer node2.Close()

	readWriter, err := p2p.ConnectNodeToStream(node2, fmt.Sprintf("/ip4/127.0.0.1/tcp/5500/p2p/%s", node1.ID()), stream.ProtocolId)
	if err != nil {
		log.Panic(err)
	}

	go stream.WriteStream(readWriter)
	go stream.ReadStream(readWriter)

	// if *targetF == "" {
	// 	attachStreamToPeer(node, protocol)


	// } else {
	// 	node.NewStream(context.Background(), )
		
	// }

	// select {}
}
