package main

import (
	"bufio"
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func createNode() host.Host {
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	return node
}

func readHelloProtocol(s network.Stream) error {
	reader := bufio.NewReader(s)
	message, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	connection := s.Conn()

	println("Message from '%s': '%s'", connection.RemotePeer().String(), message)
	return nil
}

func runTargetNode() peer.AddrInfo{	
	log.Printf("Creating target node...")
	targetNode := createNode()
	log.Printf("Target node created with ID '%s'", targetNode.ID().String())

	// TO BE IMPLEMENTED: Set stream handler for the "/hello/1.0.0" protocol
	targetNode.SetStreamHandler("/hello/1.0.0", func(s network.Stream) {
		println("Hello protocol initialized")
		err := readHelloProtocol(s)
		if err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	return *host.InfoFromHost(targetNode)
}

func runSourceNode(targetNodeInfo peer.AddrInfo) {
	log.Printf("Creating source node...")
	sourceNode := createNode()
	log.Printf("Source node created with ID '%s'", sourceNode.ID().String())

	sourceNode.Connect(context.Background(), targetNodeInfo)

	// TO BE IMPLEMENTED: Open stream and send message
	stream, err := sourceNode.NewStream(context.Background(), targetNodeInfo.ID, "/hello/1.0.0")
	if err != nil {
		panic(err)
	}

	println("Sending message...")
	message := "Hello, World!\n"
	n, err := stream.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	println(n)
}

func main() {
	ctx, _ := context.WithCancel(context.Background())

	info := runTargetNode()
	runSourceNode(info)

	<-ctx.Done()
}
