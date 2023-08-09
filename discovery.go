package main

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

func discoverPeers(ctx context.Context, h host.Host, rendezvous string) {

	log.Info("libp2p node created: ", h.ID().Pretty(), " ", h.Addrs())

	// Start peer discovery to find other peers
	log.Debug("Starting peer discovery...")

	go discoverDHTPeers(ctx, h, rendezvous)
	go discoverMDNSPeers(ctx, h, rendezvous)

	log.Info("Peer discovery completed.")

	select {}

}
