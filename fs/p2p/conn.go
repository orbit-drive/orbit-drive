package p2p

import (
	"bufio"
	"fmt"
	"os"

	discovery "github.com/libp2p/go-libp2p-discovery"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	"github.com/orbit-drive/orbit-drive/fs/pb"
	log "github.com/sirupsen/logrus"
)

const (
	// ProtocolID represents a header id for data stream processing between peers.
	ProtocolID string = "/od/sync/1.0.0"
)

func ReadHandler(s inet.Stream) {
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

func WriteHandler(s inet.Stream) {
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

// InitConn main entry for initialization of p2p connections.
func InitConn(port, nid string) error {
	lnode, err := NewLNode(port, nid)
	if err != nil {
		return err
	}
	log.WithField("host-id", lnode.ID()).Info("Host created")

	// Initialize kademlia distributed hash table from LNode host.
	kademliaDHT, err := libp2pdht.New(lnode.GetContext(), lnode)
	if err != nil {
		return err
	}
	if err = kademliaDHT.Bootstrap(lnode.GetContext()); err != nil {
		return err
	}

	lnode.ConnectToBootstrapNodes()

	// TODO: move routing to LNode ?
	log.Warn("Announcing to peers...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(lnode.GetContext(), routingDiscovery, nid)
	log.Info("Announcing successful")

	log.Warn("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(lnode.GetContext(), nid)
	if err != nil {
		return err
	}

	for peer := range peerChan {
		log.WithField("peer-id", peer.ID).Info("Peer discovered !!!")

		if peer.ID == lnode.ID() {
			continue
		}
		stream, err := lnode.NewStream(lnode.GetContext(), peer.ID, protocol.ID(ProtocolID))
		if err != nil {
			log.WithField("peer-id", peer.ID).Info("Peer connection failed")
			continue
		} else {
			lnode.AddStream(stream)
		}
		log.WithField("peer-id", peer.ID).Info("Peer connection successful")
	}

	return nil
}
