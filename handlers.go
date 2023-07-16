package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketUpgrader upgrades an HTTP connection to a WebSocket connection
var WebSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// List Topics Handler
func listTopicsHandler(c *gin.Context) {
	var topicsList []string
	topics.Range(func(key, value interface{}) bool {
		topicsList = append(topicsList, key.(string))
		return true
	})

	c.JSON(http.StatusOK, gin.H{"topics": topicsList})
}

// Create Topic Handler
func createTopicHandler(c *gin.Context) {
	var requestBody struct {
		TopicName string `json:"topicName"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topicName := requestBody.TopicName

	_, ok := topics.Load(topicName)
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Topic already exists"})
		return
	}

	pubSubTopic, err := pub.Join(topicName)
	if err != nil {
		log.Fatal(err)
	}

	topic := &Topic{
		PubSubTopic: pubSubTopic,
		Mutex:       sync.Mutex{},
	}

	topics.Store(topicName, topic)

	c.JSON(http.StatusCreated, gin.H{"topic": topicName})
}

// Get Topic Details Handler
func getTopicDetailsHandler(c *gin.Context) {
	topicName := c.Param("topicID")

	topic, ok := topics.Load(topicName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	// FIXME
	_ = topic

	c.JSON(http.StatusOK, gin.H{"topic": topicName})
}

// Join Topic Handler
func joinTopicHandler(c *gin.Context) {
	topicName := c.Param("topicID")

	topic, ok := topics.Load(topicName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	topicObj := topic.(*Topic)
	topicObj.Mutex.Lock()
	if topicObj.Conn != nil {
		topicObj.Conn.Close()
	}
	topicObj.Conn = conn
	topicObj.Mutex.Unlock()

	go handleClient(conn, topicObj)

	c.JSON(http.StatusOK, gin.H{"message": "Joined topic successfully"})
}

// List Peers Handler
func listPeersHandler(c *gin.Context) {
	topicName := c.Param("topicID")

	topic, ok := topics.Load(topicName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	topicObj := topic.(*Topic)
	topicObj.Mutex.Lock()
	defer topicObj.Mutex.Unlock()

	var peers []string
	for _, peer := range topicObj.PubSubTopic.ListPeers() {
		peers = append(peers, peer.String())
	}

	c.JSON(http.StatusOK, gin.H{"peers": peers})
}

// Connect to WebSocket Handler
func connectWebSocketHandler(c *gin.Context) {
	topicName := c.Param("topicID")

	topic, ok := topics.Load(topicName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	topicObj := topic.(*Topic)
	topicObj.Mutex.Lock()
	if topicObj.Conn != nil {
		topicObj.Conn.Close()
	}
	topicObj.Conn = conn
	topicObj.Mutex.Unlock()

	go handleClient(conn, topicObj)

	c.JSON(http.StatusOK, gin.H{"message": "WebSocket connection established"})
}
