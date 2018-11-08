package fs

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/sys"
	"github.com/wlwanpan/orbit-drive/fs/vtree"
)

// Hub represents a interface for the backend hub service.
type Hub struct {
	// HostAddr is the address of the backend hub service.
	HostAddr string

	// conn holds the websocket connection.
	conn *websocket.Conn
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
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			sys.Alert(err.Error())
		}
		log.Println(common.ToStr(message))
	}
}

// Stop closes the hub websocket connection to the backend hub.
func (h *Hub) Stop() {
	h.conn.Close()
}

// Push send a msg to websocket connection
func (h *Hub) Push(msg []byte) error {
	return nil
}