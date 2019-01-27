package p2p

import (
	"bufio"
	"context"
	"log"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	maddr "github.com/multiformats/go-multiaddr"
)

const (
	ProtocolID = ""

	AccountKey = "test-account-key"
)

func handleStream(stream inet.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readHandler(rw)
	go writeHandler(rw)
}

func readHandler(rw *bufio.ReadWriter) {}

func writeHandler(rw *bufio.ReadWriter) {}

func createHost(ctx context.Context) (host.Host, error) {
	listenMAddr, _ := maddr.NewMultiaddr("/ip4/127.0.0.1/tcp/6666")
	hostOption := libp2p.ListenAddrs(listenMAddr)

	host, err := libp2p.New(ctx, hostOption)
	if err != nil {
		return host, err
	}

	log.Printf("Host initialized: (%s)", host.ID())
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
		log.Printf("Attempting connection to: %s", peerAddr.String())
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	log.Println("Announcing...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, AccountKey)
	log.Println("Announcing successful")

	log.Println("Searching for other peers...")
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
			log.Printf("Connection to %s failed: %s", peer.ID, err.Error())
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeHandler(rw)
			go readHandler(rw)
		}

		log.Printf("Connected to %s", peer.ID)
	}

	return nil
}
