package p2p

import (
	"context"
	"fmt"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/orbit-drive/orbit-drive/pb"
	log "github.com/sirupsen/logrus"
)

// LNode represents a local node connection and is a wrapper around libp2p Host.
type LNode struct {
	host.Host
	RPC

	// Port for tcp connection of local node to listen to.
	Port string

	// NID (Network ID) is the rendez vous point for other nodes.
	NID string

	// Peers is the list of connected peers under the same NID.
	Peers []peer.ID

	// Ctx context for cancellation signal ?
	ctx context.Context
}

func NewLNode(port, nid string) *LNode {
	lnode := &LNode{
		Port:  port,
		NID:   nid,
		Peers: []peer.ID{},
		ctx:   context.Background(),
	}
	lnode.RPC = *NewRpc(lnode)
	return lnode
}

// GetPeerID returns the current local node peer id.
func (ln *LNode) GetPeerID() peer.ID {
	return ln.Host.ID()
}

// GetContext return current local node context.
func (ln *LNode) GetContext() context.Context {
	return ln.ctx
}

// AddPeer adds a new peer id to the list of connected peer id.
func (ln *LNode) AddPeer(pid peer.ID) {
	ln.Peers = append(ln.Peers, pid)
}

func (ln *LNode) initHost() error {
	addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", ln.Port)
	listenMAddr, _ := maddr.NewMultiaddr(addr)
	hostOption := libp2p.ListenAddrs(listenMAddr)

	host, err := libp2p.New(ln.GetContext(), hostOption)
	if err != nil {
		return err
	}

	ln.Host = host
	ln.initHandlers()
	return nil
}

// Request send a rpc call to connected peers.
func (ln *LNode) Request(method string) {
	var responses []*pb.Response
	var wg sync.WaitGroup

	for _, peerID := range ln.Peers {
		log.WithFields(log.Fields{
			"peer-id": peerID,
			"method":  method,
		}).Info("Request sent")

		wg.Add(1)
		go func(pid peer.ID) {
			respPb, err := ln.RequestToPeer(pid, method)
			if err != nil {
				log.Warn(err)
				return
			}
			responses = append(responses, respPb)
		}(peerID)
	}

	wg.Wait()
	log.Println(responses)
}
