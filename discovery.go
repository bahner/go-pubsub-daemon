package main

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func discoverPeers(ctx context.Context, h host.Host, rendezvous string) {

	defer wg.Done()

	log.Info("libp2p node created: ", h.ID().Pretty(), " ", h.Addrs())

	// Start peer discovery to find other peers
	log.Debug("Starting peer discovery...")

	wg.Add(2)
	go discoverDHTPeers(ctx, h, rendezvous)
	go discoverMDNSPeers(ctx, h, "")
	wg.Wait()

	log.Info("Peer discovery completed.")

	select {}

}
