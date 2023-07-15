package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"go.deanishe.net/env"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
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
)

var (
	logLevel = env.Get("MYSPACE_PUBSUB_DAEMON_LOG_LEVEL", "info")
	port     = flag.String("port", env.Get("MYSPACE_PUBSUB_DAEMON_PORT", defaultPort), "Port to listen on")
	addr     = flag.String("addr", env.Get("MYSPACE_PUBSUB_DAEMON_ADDR", defaultAddr), "Address to listen on")
)

type Topic struct {
	Mutex       sync.Mutex
	PubSubTopic *pubsub.Topic
	Conn        *websocket.Conn
}

func main() {

	ctx := context.Background()

	flag.Parse()

	listenSocket := fmt.Sprintf("%s:%s", *addr, *port)
	fmt.Println("Listening on: ", listenSocket)

	// Set log level
	lvl, err := logging.LevelFromString(logLevel)
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)
	logger := logging.Logger("myspace")
	logger.Info("Starting myspace libp2p pubsub server...")

	// Initialize libp2p host with DHT routing
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)
	if err != nil {
		log.Fatal(err)
	}

	go discoverPeers(ctx, host)

	pub, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Fatal(err)
	}

	// Create new gin engine
	router := gin.Default()

	router.GET("/topic/:topicName", func(c *gin.Context) {
		topicName := c.Param("topicName")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		t, ok := topics.Load(topicName)
		if !ok {
			pubSubTopic, err := pub.Join(topicName)
			if err != nil {
				log.Fatal(err)
			}
			t = &Topic{PubSubTopic: pubSubTopic}
			topics.Store(topicName, t)
		}

		topic := t.(*Topic)

		go handlePubSub(conn, topic)
	})

	router.Run(listenSocket)
}

func handlePubSub(conn *websocket.Conn, topic *Topic) {
	topic.Mutex.Lock()
	if topic.Conn != nil {
		topic.Conn.Close()
	}
	topic.Conn = conn
	topic.Mutex.Unlock()

	go handleClient(conn, topic)
}
func handleClient(conn *websocket.Conn, topic *Topic) {
	msg := fmt.Sprintf("Welcome to the chat room %q!", topic.PubSubTopic)
	sendText(conn, msg)

	sub, err := topic.PubSubTopic.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for reading from the WebSocket
	go func() {
		defer wg.Done()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Log the error and return from the goroutine
				log.Printf("read error: %v", err)
				return
			}

			// Publish the message to the pubsub topic
			err = topic.PubSubTopic.Publish(context.Background(), message)
			if err != nil {
				log.Printf("publish error: %v", err)
				return
			}
		}
	}()

	// Goroutine for writing to the WebSocket
	go func() {
		defer wg.Done()
		for {
			msg, err := sub.Next(context.Background())
			if err != nil {
				// Log the error and return from the goroutine
				log.Printf("subscription error: %v", err)
				return
			}

			// Write the message back to the WebSocket
			err = conn.WriteMessage(websocket.TextMessage, msg.GetData())
			if err != nil {
				log.Printf("write error: %v", err)
				return
			}
		}
	}()

	wg.Wait()
}

func sendText(c *websocket.Conn, text string) error {
	return c.WriteMessage(websocket.TextMessage, []byte(text))
}

func initDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		panic(err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT
}

func discoverPeers(ctx context.Context, h host.Host) error {
	kademliaDHT := initDHT(ctx, h)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, rendezvousString)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
		if err != nil {
			return err
		}
		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // No self connection
			}
			err := h.Connect(ctx, peer)
			if err != nil {
				fmt.Println("Failed connecting to ", peer.ID.Pretty(), ", error:", err)
			} else {
				fmt.Println("Connected to:", peer.ID.Pretty())
				anyConnected = true
			}
		}
	}
	fmt.Println("Peer discovery complete")

	return nil
}
