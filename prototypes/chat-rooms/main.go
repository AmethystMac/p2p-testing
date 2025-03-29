package main

import (
	"context"
	"flag"
	"fmt"
	"log"
)

func main() {

	ctx := context.Background()

	fmt.Print("\n\nWELCOME TO P2P CHAT APPLICATION 2.0\n\n")
	
	portF := flag.Int("p", 0, "port number")
	chatRoomNameF := flag.String("r", "general", "chat room name")

	flag.Parse()

	if *portF == 0 {
		log.Println("Error: No port number passed.")
		return
	}

	if *chatRoomNameF == "general" {
		log.Println("No chat room name argument passed.\nJoining General.")
	}
	
	node, err := CreateNode(*portF)
	if err != nil {
		log.Panicln(err)
	}
	defer node.Close()

	var ps pubsub
	ps.nodeId = node.ID()
	ps.ctx = ctx

	go ConnectToBootstrapNodes(ctx, node, *chatRoomNameF)

	topic := ps.CreateChatRoom(ctx, node, *chatRoomNameF)

	go ps.PublishToTopic(topic)
	go ps.ReadFromTopic(topic)

	// keep the terminal open
	select {}
}
