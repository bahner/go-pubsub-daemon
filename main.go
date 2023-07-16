package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.deanishe.net/env"
)

const (
	rendezvousString = "myspace"
	apiVersion       = "v0"
	multiAddr        = "/ip4/0.0.0.0/tcp/0"
)

var (
	defaultPort     = env.Get("MYSPACE_PUBSUB_DAEMON_PORT", "5002")
	defaultAddr     = env.Get("MYSPACE_PUBSUB_DAEMON_ADDR", "127.0.0.1")
	defaultLogLevel = env.Get("MYSPACE_PUBSUB_DAEMON_LOG_LEVEL", "error")
)

var (

	// API
	apiPath = fmt.Sprintf("/api/%s", apiVersion)

	// Logging for libp2p
	logger = logging.Logger("myspace")

	// Flags
	port     = flag.String("port", defaultPort, "Port to listen on")
	addr     = flag.String("addr", defaultAddr, "Address to listen on")
	logLevel = flag.String("loglevel", defaultLogLevel, "Log level for libp2p")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	// Set log level
	logging.SetLogLevel("myspace", *logLevel)
	logger.Info("Starting myspace libp2p pubsub server...")

	// Start libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(multiAddr),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start peer discovery to find other peers
	go discoverPeers(ctx, host, rendezvousString)

	// Start pubsub service
	pub, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Fatal(err)
	}

	// Create new gin engine
	router := gin.Default()

	// API Endpoints
	router.GET(apiPath+"/topics", listTopicsHandler)
	router.GET(apiPath+"/topics/:topicID/peers", listPeersHandler)
	router.GET(apiPath+"/topics/:topicID", joinTopicHandler)

	// Start server on the configured socket
	listenSocket := fmt.Sprintf("%s:%s", *addr, *port)
	router.Run(listenSocket)
}
