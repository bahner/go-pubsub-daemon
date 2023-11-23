package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/bahner/go-ma"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

const (
	apiVersion = "v0"
)

var (
	defaultRendezvous = env.Get("GO_PUBSUB_DAEMON_RENDEZVOUS", ma.RENDEZVOUS)
	defaultPort       = env.Get("GO_PUBSUB_DAEMON_PORT", "5002")
	defaultAddr       = env.Get("GO_PUBSUB_DAEMON_ADDR", "127.0.0.1")
	defaultLogLevel   = env.Get("GO_PUBSUB_DAEMON_LOG_LEVEL", "error")
)

var (

	// API
	apiPath = fmt.Sprintf("/api/%s", apiVersion)

	// Flags
	port       = flag.String("port", defaultPort, "Port to listen on")
	addr       = flag.String("addr", defaultAddr, "Address to listen on")
	logLevel   = flag.String("loglevel", defaultLogLevel, "Log level for libp2p")
	rendezvous = flag.String("rendezvous", defaultRendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	// Set log level
	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)
	log.Info("Starting go-pubsub-daemon")

	// Start libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start peer discovery to find other peers
	go discoverPeers(ctx, host, *rendezvous)

	// Start pubsub service
	pub, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Fatal(err)
	}

	// Create new gin engine
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.SetTrustedProxies(nil)
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gin.ErrorLogger())

	listenSocket := fmt.Sprintf("%s:%s", *addr, *port)

	// API Endpoints
	router.GET(apiPath+"/topics", listTopicsHandler)
	router.GET(apiPath+"/topics/:topicID/peers", listPeersHandler)
	router.GET(apiPath+"/topics/:topicID", joinTopicHandler)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = false
	config.AllowHeaders = []string{
		"Accept-Encoding:",
		"Accept-Language:",
		"Authorization",
		"Content-Length",
		"Content-Type",
		"Origin",
		"Sec-GPC:",
		"Sec-WebSocket-Extensions",
		"Sec-WebSocket-Key",
		"Sec-WebSocket-Protocol",
		"Sec-WebSocket-Version",
		"Upgrade",
		"User-Agent:",
		"X-CSRF-Token",
	}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowWildcard = true
	router.Use(cors.New(config))

	// Start server on the configured socket
	router.Run(listenSocket)
}
