package p2p

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	discovery "github.com/libp2p/go-libp2p-discovery"
	host "github.com/libp2p/go-libp2p-host"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	maddr "github.com/multiformats/go-multiaddr"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	"github.com/orbit-drive/orbit-drive/fs/pb"
	log "github.com/sirupsen/logrus"
)

const (
	// ProtocolID represents a header id for data stream processing between peers.
	ProtocolID string = "/od/sync/1.0.0"
)

func handleStream(s inet.Stream) {
	go readHandler(s)
	go writeHandler(s)
}

func readHandler(s inet.Stream) {
	for {
		data := &pb.MessageData{}
		decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
		err := decoder.Decode(data)
		if err != nil {
			log.Warn(err)
			return
		}
		log.WithFields(log.Fields{
			"msg": string(data.GetMessage()),
		}).Info("Received message from peer")
	}
}

func writeHandler(s inet.Stream) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading from stdin")
		}

		msgData := &pb.MessageData{
			Message: sendData,
		}

		writer := bufio.NewWriter(s)
		enc := protobufCodec.Multicodec(nil).Encoder(writer)
		if err = enc.Encode(msgData); err != nil {
			log.Warn(err)
			return
		}
		writer.Flush()
	}
}

func createHost(ctx context.Context, port string) (host.Host, error) {
	addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", port)
	listenMAddr, _ := maddr.NewMultiaddr(addr)
	hostOption := libp2p.ListenAddrs(listenMAddr)

	host, err := libp2p.New(ctx, hostOption)
	if err != nil {
		return host, err
	}

	host.SetStreamHandler(protocol.ID(ProtocolID), handleStream)
	return host, nil
}

func InitConn(port, nid string) error {
	ctx := context.Background()
	host, err := createHost(ctx, port)
	if err != nil {
		return err
	}
	log.WithField("host-id", host.ID()).Info("Host created")

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
		wg.Add(1)
		go func(peerAddr maddr.Multiaddr) {
			defer wg.Done()

			if err = host.Connect(ctx, *peerinfo); err != nil {
				log.WithFields(log.Fields{
					"peer-id": peerinfo.ID,
					"err-msg": err.Error(),
				}).Warn("Connection to peer failed")
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
	discovery.Advertise(ctx, routingDiscovery, nid)
	log.Info("Announcing successful")

	log.Warn("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, nid)
	if err != nil {
		return err
	}

	for peer := range peerChan {
		log.WithField("peer-id", peer.ID).Info("Peer discovered !!!")

		if peer.ID == host.ID() {
			continue
		}
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(ProtocolID))
		if err != nil {
			log.WithField("peer-id", peer.ID).Info("Peer connection failed")
			continue
		} else {
			go writeHandler(stream)
			go readHandler(stream)
		}
		log.WithField("peer-id", peer.ID).Info("Peer connection successful")
	}

	return nil
}
