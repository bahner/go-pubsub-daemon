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
	pub      *pubsub.PubSub
	topics   sync.Map
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	port     = flag.String("port", defaultPort, "Port to listen on")
	addr     = flag.String("addr", defaultAddr, "Address to listen on")
	logLevel = flag.String("loglevel", defaultLogLevel, "Log level")
	apiPath  = fmt.Sprintf("/api/%s", apiVersion)
)

func main() {
	ctx := context.Background()
	flag.Parse()

	// Set log level
	logging.SetLogLevel("myspace", *logLevel)
	logger := logging.Logger("myspace")
	logger.Info("Starting myspace libp2p pubsub server...")

	// Initialize libp2p host with DHT routing
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(multiAddr),
	)
	if err != nil {
		log.Fatal(err)
	}

	// FIMXE: go discoverPeers(ctx, host, rendezvousString)

	pub, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Fatal(err)
	}

	// Create new gin engine
	router := gin.Default()

	// API Endpoints
	router.GET(apiPath+"/topics", listTopicsHandler)
	// FIXME: router.POST(apiPath+"/topics/:topicID", publishMessageHandler)
	router.GET(apiPath+"/topics/:topicID/peers", listPeersHandler)
	router.GET(apiPath+"/topics/:topicID", joinTopicHandler)

	listenSocket := fmt.Sprintf("%s:%s", *addr, *port)
	router.Run(listenSocket)
}
