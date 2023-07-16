package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

// discoverPeers performs peer discovery using the DHT and connects to discovered peers
func discoverPeers(ctx context.Context, h host.Host, rendezvousString string) error {
	dhtInstance, err := initDHT(ctx, h)
	if err != nil {
		return err
	}

	routingDiscovery := drouting.NewRoutingDiscovery(dhtInstance)
	dutil.Advertise(ctx, routingDiscovery, rendezvousString)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
		if err != nil {
			return fmt.Errorf("peer discovery error: %w", err)
		}

		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // Skip self connection
			}

			err := h.Connect(ctx, peer)
			if err != nil {
				fmt.Printf("Failed connecting to %s, error: %v\n", peer.ID.Pretty(), err)
			} else {
				fmt.Println("Connected to:", peer.ID.Pretty())
				anyConnected = true
			}
		}
	}
	fmt.Println("Peer discovery complete")

	return nil
}
