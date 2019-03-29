package p2p

import (
	"bufio"

	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	"github.com/orbit-drive/orbit-drive/fs/pb"
	"github.com/orbit-drive/orbit-drive/fsutil"
	log "github.com/sirupsen/logrus"
)

const (
	// ProtocolRequestID - protocol header id for request traffic
	ProtocolRequestID string = "/od/syncreq/1.0.0"

	// ProtocolResponseID - protocol header id for response traffic
	ProtocolResponseID string = "/od/syncresp/1.0.0"
)

type ReqResp struct {
	requestPb *pb.Request
	respChan  chan *pb.Response
}

func newReqResp(reqPb *pb.Request) *ReqResp {
	return &ReqResp{
		requestPb: reqPb,
		respChan:  make(chan *pb.Response),
	}
}

// RPC handles the data stream from/to libp2p host and stores all
// incoming/outgoing request until fully resolve or timeout is reached.
type RPC struct {
	lnode *LNode

	// TODO: Might need to create a self contained context queue to timeout request pending for x amount of seconds.
	// ReqOut represents the queue outgoing requests from the current node.
	ReqOut map[string]*ReqResp

	// ReqIn represents a chan of incoming request to process.
	ReqIn chan *pb.Request
}

func NewRpc(lnode *LNode) *RPC {
	return &RPC{
		lnode:  lnode,
		ReqOut: make(map[string]*ReqResp),
		ReqIn:  make(chan *pb.Request),
	}
}

// RequestID returns the procotol resquest id.
func (rpc *RPC) RequestID() protocol.ID {
	return protocol.ID(ProtocolRequestID)
}

// ResponseID returns the protol response id.
func (rpc *RPC) ResponseID() protocol.ID {
	return protocol.ID(ProtocolResponseID)
}

func (rpc *RPC) initHandlers() {
	rpc.lnode.SetStreamHandler(rpc.RequestID(), rpc.reqHandler)
	rpc.lnode.SetStreamHandler(rpc.ResponseID(), rpc.respHandler)
}

func (rpc *RPC) createReq(method string) *pb.Request {
	return &pb.Request{
		PeerId:    string(rpc.lnode.GetPeerID()),
		RequestId: fsutil.RandUUID(),
		Method:    method,
	}
}

func (rpc *RPC) registerReqOut(reqPb *pb.Request) *ReqResp {
	reqResp := newReqResp(reqPb)
	rpc.ReqOut[reqPb.GetRequestId()] = reqResp
	return reqResp
}

func (rpc *RPC) findReqOut(reqID string) *ReqResp {
	reqResp, ok := rpc.ReqOut[reqID]
	if ok {
		return reqResp
	}
	return nil
}

// RequestToPeer opens a stream to a single peer and sends a proto request.
func (rpc *RPC) RequestToPeer(peerID peer.ID, method string) (*pb.Response, error) {
	stream, err := rpc.lnode.NewStream(rpc.lnode.GetContext(), peerID, rpc.RequestID())
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(stream)
	requestPayload := rpc.createReq(method)

	enc := protobufCodec.Multicodec(nil).Encoder(writer)
	if err = enc.Encode(requestPayload); err != nil {
		return nil, err
	}
	writer.Flush()

	reqResp := rpc.registerReqOut(requestPayload)
	respPb := <-reqResp.respChan
	return respPb, nil
}

// reqHandler: remote peer request handler (received request from peer)
func (rpc *RPC) reqHandler(s inet.Stream) {
	req := &pb.Request{}
	reader := bufio.NewReader(s)
	decoder := protobufCodec.Multicodec(nil).Decoder(reader)
	if err := decoder.Decode(req); err != nil {
		log.Warn(err)
		return
	}
	log.WithFields(log.Fields{
		"peer-id": req.GetPeerId(),
		"req-id":  req.GetRequestId(),
		"method":  req.GetMethod(),
	}).Info("Received request from peer")
}

func (rpc *RPC) respHandler(s inet.Stream) {
	resp := &pb.Response{}
	reader := bufio.NewReader(s)
	decoder := protobufCodec.Multicodec(nil).Decoder(reader)
	if err := decoder.Decode(resp); err != nil {
		log.Warn(err)
		return
	}

	reqResp := rpc.findReqOut(resp.GetRequestId())
	if reqResp != nil {
		reqResp.respChan <- resp
		return
	}

	log.WithFields(log.Fields{
		"peer-id":      resp.GetPeerId(),
		"request-uuid": resp.GetRequestId(),
	}).Warn("Received response from peer with no corresponding request")
}
