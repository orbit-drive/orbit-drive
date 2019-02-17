package p2p

import (
	"bufio"
	"context"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	maddr "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const (
	ProtocolID = "/od/sync/1.0.0"

	AccountKey = "test-account-key"
)

func handleStream(stream inet.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readHandler(rw)
	go writeHandler(rw)
}

func readHandler(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Warn(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			log.WithFields(log.Fields{
				"msg": str,
			}).Info("Received message from peer")
		}
	}
}

func writeHandler(rw *bufio.ReadWriter) {
	for {
		_, err := rw.WriteString("this is a test \n")
		if err != nil {
			log.Warn(err)
		}

		// Why flush buffer ??
		if err = rw.Flush(); err != nil {
			log.Warn(err)
		}
	}
}

func createHost(ctx context.Context) (host.Host, error) {
	listenMAddr, _ := maddr.NewMultiaddr("/ip4/127.0.0.1/tcp/6666")
	hostOption := libp2p.ListenAddrs(listenMAddr)

	host, err := libp2p.New(ctx, hostOption)
	if err != nil {
		return host, err
	}

	// log.WithField("host-id", host.ID()).Info("Host created")
	host.SetStreamHandler(protocol.ID(ProtocolID), handleStream)
	return host, nil
}

func InitConn() error {
	ctx := context.Background()
	host, err := createHost(ctx)
	if err != nil {
		return err
	}

	kademliaDHT, err := libp2pdht.New(ctx, host)
	if err != nil {
		return err
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range getBootstrapAddrs() {
		peerinfo, _ := peerstore.InfoFromP2pAddr(peerAddr)
		// log.WithField("peer-addr", peerAddr.String()).Warning("Attempting connection to peer")
		wg.Add(1)
		go func(peerAddr maddr.Multiaddr) {
			defer wg.Done()
			if err = host.Connect(ctx, *peerinfo); err != nil {
				log.Warn(err)
			} else {
				log.WithFields(log.Fields{
					"peer-id":   peerinfo.ID,
					"peer-addr": peerAddr.String(),
				}).Info("Connection established with peer")
			}
		}(peerAddr)
	}
	wg.Wait()

	log.Warn("Announcing to peers...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, AccountKey)
	log.Info("Announcing successful")

	log.Warn("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, AccountKey)
	if err != nil {
		return err
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(ProtocolID))

		if err != nil {
			log.WithFields(log.Fields{
				"peer-id": peer.ID,
				"err-msg": err.Error(),
			}).Warn("Connection to peer failed")
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeHandler(rw)
			go readHandler(rw)
		}

		log.WithField("peer-id", peer.ID).Info("Peer connection successful")
	}

	return nil
}
