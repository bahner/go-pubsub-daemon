package main

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func initMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}

func discoverMDNSPeers(ctx context.Context, h host.Host, rendezvous string) chan peer.AddrInfo {
	anyConnected := false
	for !anyConnected {

		log.Info("Starting MDNS peer discovery.")
		peerChan := initMDNS(h, rendezvous)

		for peer := range peerChan {
			log.Debugf("Found peer: %s\n", peer.ID.String())
			if peer.ID == h.ID() {
				continue // Skip self connection
			}

			err := h.Connect(ctx, peer)
			if err != nil {
				log.Debugf("Failed connecting to MDNS peer %s, error: %v\n", peer.ID.String(), err)
			} else {
				log.Infof("Connected to MDNS peer: %s", peer.ID.String())
				anyConnected = true
			}
		}
	}

	log.Info("MDNS Peer discovery complete")

	return nil

}
