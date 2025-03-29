package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
)

// Create and initialize the Kademlia DHT
func InitializeDHT(ctx context.Context, node host.Host) *dht.IpfsDHT {

	// Creating a custom kademliaDHT DHT
	kademliaDHT, err := dht.New(ctx, node)
	if err != nil {
		log.Panicln(err)
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Panicln(err)
	}

	// Connecting to default bootstrap nodes
	var wg sync.WaitGroup
	for _, multiAddr := range dht.DefaultBootstrapPeers {
		peer, err := peer.AddrInfoFromP2pAddr(multiAddr)
		if err != nil {
			log.Panicln(err)
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			if err = node.Connect(ctx, *peer); err != nil {
				log.Panicln(err)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT
}

// Connect the [node] to the chat room
func ConnectToBootstrapNodes(ctx context.Context, node host.Host, chatRoomName string) {

	kademliaDHT := InitializeDHT(ctx, node)

	// Announce ourselves and find peers using bootstrap peers to join the chat room
	routeDiscovery := routing.NewRoutingDiscovery(kademliaDHT)
	util.Advertise(ctx, routeDiscovery, chatRoomName)

	peersFound := false
	for !peersFound {
		roomPeers, err := routeDiscovery.FindPeers(ctx, chatRoomName)
		if err != nil {
			log.Panicln(err)
		}
		for peer := range roomPeers {
			if peer.ID == node.ID() {
				continue
			}
			err := node.Connect(ctx, peer)
			if err != nil {
				fmt.Printf("Node cannot connect to %s\n", peer.ID)
			} else {
				fmt.Printf("Node connected to %s\n", peer.ID)
				peersFound = true
			}
		}
	}
	fmt.Println("Peer discovery complete")
}
