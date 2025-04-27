package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
)

func main() {

	ctx := context.Background()

	fmt.Print("\n\nWELCOME TO P2P CHAT APPLICATION 2.2\n\n")
	
	HandleFlags()

	port, _ := strconv.Atoi(GetConfig("port"))
	chatRoomName := GetConfig("chatRoomName")

	if port == 0 {
		log.Println("Error: No port number passed.")
		return
	}

	if chatRoomName == "general" {
		log.Println("No chat room name argument passed.\nJoining General.")
	}
	
	node, err := CreateNode(port)
	if err != nil {
		log.Panicln(err)
	}
	defer node.Close()

	var ps Pubsub
	ps.nodeId = node.ID()
	ps.ctx = ctx

	go ConnectToBootstrapNodes(ctx, node, chatRoomName)

	topic := ps.CreateChatRoom(ctx, node, chatRoomName)

	go ps.PublishToTopic(topic)
	go ps.ReadFromTopic(topic)

	// keep the terminal open
	select {}
}
