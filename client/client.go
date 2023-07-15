package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("baseurl", "localhost:8080", "http service address")
var topic = flag.String("topicname", "myspace", "websocket topic")

func main() {
	flag.Parse()

	// The WebSocket URL
	u := fmt.Sprintf("ws://%s/topic/%s", *addr, *topic)

	// Connect to the WebSocket server
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) == "quit" {
			break
		}

		err := c.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("error:", err)
	}
}
