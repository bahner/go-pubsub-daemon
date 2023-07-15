package main

import (
	"context"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	libp2p "github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var pub *pubsub.PubSub
var topics sync.Map

func main() {
	// Initialize libp2p pubsub
	host, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}

	pub, err = pubsub.NewGossipSub(context.Background(), host)
	if err != nil {
		log.Fatal(err)
	}

	// Create new gin engine
	router := gin.Default()

	// Handle /topic/:topicName
	router.GET("/topic/:topicName", func(c *gin.Context) {
		topicName := c.Param("topicName")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		t, ok := topics.Load(topicName)
		if !ok {
			t, err = pub.Join(topicName)
			if err != nil {
				log.Fatal(err)
			}
			topics.Store(topicName, t)
		}

		topic := t.(*pubsub.Topic)

		go handlePubSub(conn, topic)
	})

	router.Run()
}

func handlePubSub(conn *websocket.Conn, topic *pubsub.Topic) {
	sub, err := topic.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Cancel()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// publish the message
		err = topic.Publish(context.Background(), message)
		if err != nil {
			log.Fatal(err)
		}

		// Receive the next message from the pubsub
		msg, err := sub.Next(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// Send the message over the websocket connection
		if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
			return
		}
	}
}
