package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var (
	pub      *pubsub.PubSub
	topics   sync.Map
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// List Topics Handler
func listTopicsHandler(c *gin.Context) {
	var topicsList []string
	topics.Range(func(key, value interface{}) bool {
		topicsList = append(topicsList, key.(string))
		return true
	})

	c.JSON(http.StatusOK, gin.H{"topics": topicsList})
}

func getOrCreateTopic(topicID string) (*Topic, error) {
	topic, ok := topics.Load(topicID)
	if ok {
		if t, ok := topic.(*Topic); ok {
			return t, nil
		}
	}

	pubSubTopic, err := pub.Join(topicID)
	if err != nil {
		return nil, err
	}

	topic = &Topic{
		PubSubTopic: pubSubTopic,
		Mutex:       sync.Mutex{},
		TopicID:     topicID,
	}

	topics.Store(topicID, topic)

	return topic.(*Topic), nil
}

// Join Topic Handler
func joinTopicHandler(c *gin.Context) {
	topicID := c.Param("topicID")

	topic, err := getOrCreateTopic(topicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acquiring topic failed"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	topic.Mutex.Lock()
	topic.Conn = conn // Assign the new connection to the topic
	topic.Mutex.Unlock()

	go handleClient(conn, topic) // Move handleClient inside the critical section
}

// List Peers Handler
func listPeersHandler(c *gin.Context) {
	topicID := c.Param("topicID")

	topic, ok := topics.Load(topicID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	topicObj := topic.(*Topic)
	topicObj.Mutex.Lock()
	defer topicObj.Mutex.Unlock()

	var peers []string
	for _, peer := range topicObj.PubSubTopic.ListPeers() {
		fmt.Printf("TopicPeer: %s\n", peer)
		peers = append(peers, peer.String())
	}

	c.JSON(http.StatusOK, gin.H{"peers": peers})
}
