package p2p

import (
	"context"
	"fmt"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	maddr "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

// LNode represents a local node connection and is a wrapper around libp2p Host.
type LNode struct {
	host.Host

	// Port for tcp connection of local node to listen to.
	Port string

	// NID (Network ID) is the rendez vous point for other nodes.
	NID string

	// Streams holds a list of streams from other nodes.
	Streams []inet.Stream

	// Ctx context for cancellation signal ?
	ctx context.Context
}

func NewLNode(port, nid string) (*LNode, error) {
	lnode := &LNode{
		Port:    port,
		NID:     nid,
		Streams: []inet.Stream{},
		ctx:     context.Background(),
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
	protocolID := protocol.ID(ProtocolID)

	host, err := libp2p.New(ln.GetContext(), hostOption)
	if err != nil {
		return err
	}

	host.SetStreamHandler(protocolID, ln.AddStream)
	ln.Host = host
	return nil
}

func (ln *LNode) AddStream(s inet.Stream) {
	ln.Streams = append(ln.Streams, s)

	// TODO: replace with actual proto handlers.
	go ReadHandler(s)
	go WriteHandler(s)
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
