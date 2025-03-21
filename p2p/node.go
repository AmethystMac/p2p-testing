package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"blockchain-prototype/stream"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

func CreateNode(port int) (host.Host, error) {

	private, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	opts := []libp2p.Option {
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", port)),
		libp2p.Identity(private),
	}

	node, err := libp2p.New(opts...)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	return node, err
}


func CreateNewStreamWithNode(node host.Host, protocol protocol.ID) {
	node.SetStreamHandler("/chat/1.0.0", stream.ChatProtocolHandler)
}

func ConnectNodeToStream(node host.Host, destination string, protocol protocol.ID) (*bufio.ReadWriter, error) {

	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		log.Panic(err)
	}
	
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	node.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	
	streamRef, err := node.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	readWriter := bufio.NewReadWriter(bufio.NewReader(streamRef), bufio.NewWriter(streamRef))

	return readWriter, nil
}
