package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

// Creates a node using [port]
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

	fmt.Printf("\nNode has been created with ID: %s", node.ID())

	return node, err
}
