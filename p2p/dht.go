package p2p

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	sys "blockchain-prototype/sys"
	utils "blockchain-prototype/utils"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

// Setting the bootstrap peer addrs
func GetBootstrapPeerAddrs() []peer.AddrInfo {

	var BootstrapPeers []peer.AddrInfo;
	for _, s := range []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
		"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			log.Panicln(err)
		}
		info, err := peer.AddrInfoFromP2pAddr(ma)
		if err != nil {
			log.Panicln(err)
		}
		BootstrapPeers = append(BootstrapPeers, *info)
	}

	return BootstrapPeers
}

// Create and initialize the Kademlia DHT
func InitializeDHT(ctx context.Context, node host.Host) *dht.IpfsDHT {

	opts := []dht.Option {
		// dht.BootstrapPeers(GetBootstrapPeerAddrs()...),
		// dht.ProtocolPrefix("/amethyst"),
		// dht.BucketSize(25),
	}

	// Creating a custom kademliaDHT DHT
	kademliaDHT, err := dht.New(ctx, node, opts...)
	if err != nil {
		log.Panicln(err)
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Panicln(err)
	}

	// Connecting to bootstrap nodes
	// go kademliaDHT.Bootstrap(ctx)
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

// Connect the [node] to the chatroom
func ConnectToBootstrapNodes(ctx context.Context, node host.Host, chatRoomName string) {

	kademliaDHT := InitializeDHT(ctx, node)

	// Announce ourselves and find peers using bootstrap peers to join the chatroom
	routeDiscovery := routing.NewRoutingDiscovery(kademliaDHT)
	util.Advertise(ctx, routeDiscovery, chatRoomName)

	peersFound := false
	logs, _ := strconv.ParseBool(utils.GetConfig("logs"))

	sys.SetProcessStatus("dht", "started")
	fmt.Println("Searching for peers...")
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
				if logs {
					fmt.Printf("Node cannot connect to %s\n", peer.ID)
				}
			} else {
				fmt.Printf("Node connected to %s\n", peer.ID)
				peersFound = true
			}
		}
	}
	fmt.Println("Peer discovery complete.")
	sys.SetProcessStatus("dht", "done")
}
