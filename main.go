package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	p2p "blockchain-prototype/p2p"
	utils "blockchain-prototype/utils"
)

func main() {

	ctx := context.Background()

	fmt.Print("\n\nWELCOME TO P2P CHAT APPLICATION 2.2\n\n")
	
	utils.HandleFlags()

	port, _ := strconv.Atoi(utils.GetConfig("port"))
	chatRoomName := utils.GetConfig("chatRoomName")

	if port == 0 {
		log.Println("Error: No port number passed.")
		return
	}

	if chatRoomName == "general" {
		log.Println("No chat room name argument passed.\nJoining General.")
	}
	
	node, err := p2p.CreateNode(port)
	if err != nil {
		log.Panicln(err)
	}
	defer node.Close()

	var ps p2p.Pubsub
	ps.NodeID = node.ID()
	ps.Ctx = ctx

	go p2p.ConnectToBootstrapNodes(ctx, node, chatRoomName)

	topic := ps.CreateChatRoom(ctx, node, chatRoomName)

	go ps.PublishToTopic(topic)
	go ps.ReadFromTopic(topic)

	// keep the terminal open
	select {}
}
