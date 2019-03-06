package p2p

import (
	"bufio"
	"context"
	"fmt"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	maddr "github.com/multiformats/go-multiaddr"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	"github.com/orbit-drive/orbit-drive/fs/pb"
	log "github.com/sirupsen/logrus"
)

const (
	// ProtocolRequestID - protocol header id for request traffic
	ProtocolRequestID string = "/od/syncreq/1.0.0"

	// ProtocolResponseID - protocol header id for response traffic
	ProtocolResponseID string = "/od/syncresp/1.0.0"
)

// LNode represents a local node connection and is a wrapper around libp2p Host.
type LNode struct {
	host.Host

	// Port for tcp connection of local node to listen to.
	Port string

	// NID (Network ID) is the rendez vous point for other nodes.
	NID string

	// Peers is the list of connected peers under the same NID.
	Peers []peer.ID

	// Ctx context for cancellation signal ?
	ctx context.Context
}

func NewLNode(port, nid string) (*LNode, error) {
	lnode := &LNode{
		Port:  port,
		NID:   nid,
		Peers: []peer.ID{},
		ctx:   context.Background(),
	}
	if err := lnode.initHost(); err != nil {
		return nil, err
	}

	return lnode, nil
}

func (ln *LNode) initHost() error {
	addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", ln.Port)
	listenMAddr, _ := maddr.NewMultiaddr(addr)
	hostOption := libp2p.ListenAddrs(listenMAddr)

	host, err := libp2p.New(ln.GetContext(), hostOption)
	if err != nil {
		return err
	}

	host.SetStreamHandler(protocol.ID(ProtocolRequestID), ln.requestHandler)
	host.SetStreamHandler(protocol.ID(ProtocolResponseID), ln.responseHandler)
	ln.Host = host
	return nil
}

// requestHandler: remote peer request handler (received request from peer)
func (ln *LNode) requestHandler(s inet.Stream) {
	data := &pb.MessageData{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	if err != nil {
		log.Warn(err)
		return
	}
	log.WithField("msg", string(data.GetMessage())).Info("Received message from peer")
}

// responseHandler: remote peer response handler (received response from peer)
func (ln *LNode) responseHandler(s inet.Stream) {}

func (ln *LNode) BroadcastRequest(msg string) {
	msgData := &pb.MessageData{
		Message: msg,
	}

	for _, peerID := range ln.Peers {
		s, err := ln.NewStream(ln.GetContext(), peerID, protocol.ID(ProtocolRequestID))
		if err != nil {
			log.Warn(err)
		}
		writer := bufio.NewWriter(s)
		enc := protobufCodec.Multicodec(nil).Encoder(writer)
		if err = enc.Encode(msgData); err != nil {
			log.Warn(err)
		}
		writer.Flush()
	}
}

func (ln *LNode) AddPeer(pid peer.ID) {
	ln.Peers = append(ln.Peers, pid)
}

func (ln *LNode) GetContext() context.Context {
	return ln.ctx
}

// ConnectToBootstrapNodes initialize connection to hardcoded ipfs nodes addr.
func (ln *LNode) ConnectToBootstrapNodes() {
	var wg sync.WaitGroup
	for _, peerAddr := range getBootstrapAddrs() {
		peerinfo, _ := peerstore.InfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func(peerAddr maddr.Multiaddr) {
			defer wg.Done()

			err := ln.Connect(ln.GetContext(), *peerinfo)
			if err != nil {
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
}
