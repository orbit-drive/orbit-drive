package p2p

import (
	discovery "github.com/libp2p/go-libp2p-discovery"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/orbit-drive/orbit-drive/sys"
	log "github.com/sirupsen/logrus"
)

var lnode *LNode

// InitConn main entry for initialization of p2p connections.
func InitConn(port, nid string) error {
	lnode = NewLNode(port, nid)
	if err := lnode.initHost(); err != nil {
		return err
	}
	log.WithField("host-id", lnode.ID()).Info("Host created")

	ctx := lnode.GetContext()
	// Initialize kademlia distributed hash table from LNode host.
	kademliaDHT, err := libp2pdht.New(ctx, lnode)
	if err != nil {
		return err
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return err
	}

	// Connect lnode to bootstrap libp2p nodes
	ConnectToBootstrapNodes(lnode)

	// TODO: move routing to LNode ?
	log.Warn("Announcing to peers...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, nid)
	log.Info("Announcing successful")

	log.Warn("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, nid)
	if err != nil {
		return err
	}

	for peer := range peerChan {
		log.WithField("peer-id", peer.ID).Info("Peer discovered!")

		if peer.ID == lnode.ID() {
			continue
		}

		lnode.AddPeer(peer.ID)
		sys.Notify("Peer connected: ", string(peer.ID))
	}

	return nil
}
