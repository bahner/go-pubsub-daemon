package main

import (
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
	if topic.Conn != nil {
		topic.Conn.Close()
	}
	topic.Conn = conn
	topic.Mutex.Unlock()

	go handleClient(conn, topic)

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

// // Webocket Handler
// func webSocketHandler(c *gin.Context) {
// 	topicName := c.Param("topicID")

// 	topic, ok := topics.Load(topicName)
// 	if !ok {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
// 		return
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
// 		return
// 	}

// 	topicObj := topic.(*Topic)
// 	topicObj.Mutex.Lock()
// 	if topicObj.Conn != nil {
// 		topicObj.Conn.Close()
// 	}
// 	topicObj.Conn = conn
// 	topicObj.Mutex.Unlock()

// 	go handleClient(conn, topicObj)

// 	c.JSON(http.StatusOK, gin.H{"message": "WebSocket connection established"})
// }

// // Publish Message Handler
// func publishMessageHandler(c *gin.Context) {
// 	var requestBody struct {
// 		TopicName string `json:"topicName"`
// 		Message   string `json:"message"`
// 	}

// 	if err := c.ShouldBindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	topicName := requestBody.TopicName
// 	message := requestBody.Message

// 	topic, ok := topics.Load(topicName)
// 	if !ok {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
// 		return
// 	}

// 	topicObj := topic.(*Topic)
// 	topicObj.Mutex.Lock()
// 	defer topicObj.Mutex.Unlock()

// 	if topicObj.Conn == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "No connections available"})
// 		return
// 	}

// 	if err := topicObj.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Message published successfully"})
// }
