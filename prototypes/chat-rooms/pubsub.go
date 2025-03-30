package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	p2p_pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// The pubsub struct allows you to interact with the GossipSub protocol
type pubsub struct {
	nodeId peer.ID;
	ctx context.Context;
}

// Creates a new chat room with [chatRoomName] topic
func (ps *pubsub) CreateChatRoom(ctx context.Context, node host.Host, chatRoomName string) *p2p_pubsub.Topic {
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

// Reads messages from the [topic] chat room
func (ps *pubsub) ReadFromTopic(topic *p2p_pubsub.Topic) {
	
	sub, err := topic.Subscribe()
	if err != nil {
		log.Panicln(err)
	}
	sub.Next(ps.ctx)

	for {
		message, err := sub.Next(ps.ctx)
		if err != nil {
			log.Panicln(err)
		}
		
		if message.ReceivedFrom != ps.nodeId {
			fmt.Printf("%s> ", message.Data)
		}
	}
}

// Writes message to the [topic] chat room
func (ps *pubsub) PublishToTopic(topic *p2p_pubsub.Topic) {

	cmdLine := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := cmdLine.ReadString('\n')
		if err != nil {
			log.Panicln(err)
		}

		if err = topic.Publish(ps.ctx, []byte(input)); err != nil {
			log.Panicln(err)
		}
	}
}
