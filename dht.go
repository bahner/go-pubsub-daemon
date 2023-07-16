package main

import (
	"context"
	"fmt"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

// initDHT initializes the DHT and bootstraps it
func initDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {
	dhtInstance, err := dht.New(ctx, h)
	if err != nil {
		return nil, fmt.Errorf("failed to create DHT instance: %w", err)
	}

	if err := dhtInstance.Bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("failed to bootstrap DHT: %w", err)
	}

	return dhtInstance, nil
}
