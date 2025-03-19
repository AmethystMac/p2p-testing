package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func chatProtocolHandler(s network.Stream) {

	reader := bufio.NewReader(s)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Panicln(err)
	}

	log.Println(message)
}

func createNode(port int) (host.Host, error) {

	private, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	opts := []libp2p.Option {
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", port)),
		libp2p.Identity(private),
	}

	node, err := libp2p.New(opts...)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	return node, err
}

func attachStreamToPeer(peer host.Host, protocol protocol.ID) {
	peer.SetStreamHandler(protocol, chatProtocolHandler)
}