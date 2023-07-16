package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	rendezvousString = "myspace"
	defaultPort      = "5002"
	defaultAddr      = "127.0.0.1"
)

var (
	pub      *pubsub.PubSub
	topics   sync.Map
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	port = flag.String("port", defaultPort, "Port to listen on")
	addr = flag.String("addr", defaultAddr, "Address to listen on")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	listenSocket := fmt.Sprintf("%s:%s", *addr, *port)
	fmt.Println("Listening on:", listenSocket)

	// Set log level
	logging.SetLogLevel("myspace", "info")
	logger := logging.Logger("myspace")
	logger.Info("Starting myspace libp2p pubsub server...")

	// Initialize libp2p host with DHT routing
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%s", *addr, *port)),
	)
	if err != nil {
		log.Fatal(err)
	}

	go discoverPeers(ctx, host, rendezvousString)

	pub, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Fatal(err)
	}

	// Create new gin engine
	router := gin.Default()

	// API Endpoints
	router.GET("/topics", listTopicsHandler)
	router.POST("/topics", createTopicHandler)
	router.GET("/topics/:topicID", getTopicDetailsHandler)
	router.GET("/topics/:topicID/peers", listPeersHandler)
	router.GET("/topics/:topicID/connect", connectWebSocketHandler)
	router.GET("/topics/:topicID/join", joinTopicHandler)

	router.Run(listenSocket)
}
