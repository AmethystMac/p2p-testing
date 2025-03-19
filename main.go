package main

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	
	// _ = flag.Int("p", 0, "port number")
	// targetF := flag.String("t", "", "target address")

	const protocol protocol.ID = "/chat/1.0.0"

	node1, err := createNode(5500)
	if err != nil {
		log.Panic(err)
	}
	defer node1.Close()

	attachStreamToPeer(node1, protocol)

	node2, err := createNode(5600)
	if err != nil {
		log.Panic(err)
	}
	defer node2.Close()

	maddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/5500/p2p/%s", node1.ID()))
	if err != nil {
		log.Panic(err)
	}
	
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("\n\n\n%s\n\n\n%s\n\n\n", info.ID, info.Addrs)

	node2.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	log.Printf("%s\n\n\n", node2.Peerstore().Peers())

	streamRef, err := node2.NewStream(context.Background(), info.ID, protocol)
	if err != nil {
		log.Panic(err)
	}

	streamRef.Write([]byte("Matthew the Great\n"))

	// if *targetF == "" {
	// 	attachStreamToPeer(node, protocol)


	// } else {
	// 	node.NewStream(context.Background(), )
		
	// }
}
