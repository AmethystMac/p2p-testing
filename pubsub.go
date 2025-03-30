package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	p2p_pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// The Pubsub struct allows you to interact with the GossipSub protocol
type Pubsub struct {
	nodeId		peer.ID;
	ctx			context.Context;
}

type Message struct {
	SenderID    string `json:"sender_id"`
	Content     string `json:"content"`
}

// Creates a new chatroom with [chatRoomName] topic
func (ps *Pubsub) CreateChatRoom(ctx context.Context, node host.Host, chatRoomName string) *p2p_pubsub.Topic {

	pubsub, err := p2p_pubsub.NewGossipSub(ctx, node)
	if err != nil {
		log.Panicln(err)
	}

	topic, err := pubsub.Join(chatRoomName)
	if err != nil {
		log.Panicln(err)
	}

	return topic
}

// Reads messages from the [topic] chatroom
func (ps *Pubsub) ReadFromTopic(topic *p2p_pubsub.Topic) {
	
	sub, err := topic.Subscribe()
	if err != nil {
		log.Panicln(err)
	}

	for {
		message, err := sub.Next(ps.ctx)
		if err != nil {
			log.Panicln(err)
		}

		var json_decode Message
		if err = json.Unmarshal(message.Data, &json_decode); err != nil {
			fmt.Println("Error: Cannot decode the message.")
		}
		
		if message.ReceivedFrom != ps.nodeId {
			fmt.Printf("%s: %s> ", json_decode.SenderID, json_decode.Content)
		}
	}
}

// Writes message to the [topic] chatroom
func (ps *Pubsub) PublishToTopic(topic *p2p_pubsub.Topic) {

	cmdLine := bufio.NewReader(os.Stdin)
	userName := GetConfig("userName")

	senderId := ps.nodeId.String()
	if userName != "" {
		senderId = userName
	}

	for {
		fmt.Print("> ")
		input, err := cmdLine.ReadString('\n')
		if err != nil {
			log.Panicln(err)
		}

		m := Message {
			SenderID: senderId,
			Content: input,
		}
	
		json_encode, err := json.Marshal(m)
		if err != nil {
			log.Panicln(err)
		}

		if err = topic.Publish(ps.ctx, []byte(json_encode)); err != nil {
			log.Panicln(err)
		}
	}
}
