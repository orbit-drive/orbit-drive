package fs

import (
	"log"
	"net/url"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/sys"
	"github.com/wlwanpan/orbit-drive/fs/vtree"
	"github.com/wlwanpan/orbit-drive/pb"
)

// Hub represents a interface for the backend hub service.
type Hub struct {
	// HostAddr is the address of the backend hub service.
	HostAddr string

	// conn holds the websocket connection.
	conn *websocket.Conn

	// updates
	updates chan []byte
}

// NewHub creates and start a websocket connection to backend hub.
func NewHub(addr string) (*Hub, error) {
	hub := &Hub{HostAddr: addr}
	if err := hub.Connect(); err != nil {
		return &Hub{}, nil
	}
	return hub, nil
}

// Connect dial the backend hub and establish a websocket connection
// and stores the connection to the hub conn.
func (h *Hub) Connect() error {
	u := url.URL{
		Scheme: "ws",
		Host:   h.HostAddr,
		Path:   "/account/subscribe",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	h.conn = conn
	return nil
}

// Sync listens to incoming traffics from the backend hub and
// call the appropriate handler to mutate the vtree.
func (h *Hub) Sync(vt *vtree.VTree) {
	for {
		_, msg, err := h.conn.ReadMessage()
		if err != nil {
			sys.Alert(err.Error())
		}
		log.Println(common.ToStr(msg))
		h.updates <- msg
	}
}

// Updates returns a parsed channel, parsing ws bytes to proto hub message.
func (h *Hub) Updates() (<-chan pb.Payload, <-chan error) {
	updates := make(chan pb.Payload)
	errs := make(chan error)
	go func() {
		update := <-h.updates
		hubMsg := &pb.Payload{}
		err := proto.Unmarshal(update, hubMsg)
		if err != nil {
			errs <- err
		}
		updates <- *hubMsg
	}()
	return updates, errs
}

// Stop closes the hub websocket connection to the backend hub.
func (h *Hub) Stop() {
	defer h.conn.Close()
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := h.conn.WriteMessage(websocket.CloseMessage, closeMsg)
	if err != nil {
		sys.Alert(err.Error())
	}
}

// Push send a msg to websocket connection
func (h *Hub) Push(msg []byte) error {
	return h.conn.WriteMessage(websocket.TextMessage, msg)
}
