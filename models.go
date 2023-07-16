package main

import (
	"sync"

	"github.com/gorilla/websocket"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Topic represents a chatroom topic
type Topic struct {
	Mutex       sync.Mutex
	PubSubTopic *pubsub.Topic
	Conn        *websocket.Conn
	Details     TopicDetails
}

type TopicDetails struct {
	TopicName string
	TopicID   string
}
