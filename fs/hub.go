package fs

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/orbit-drive/orbit-drive/common"
	"github.com/orbit-drive/orbit-drive/fs/sys"
	"github.com/orbit-drive/orbit-drive/fs/vtree"
)

var (
	// ErrHubConnFailed is returned when connection to hub is not established.
	ErrHubConnFailed = errors.New("hub connection failed")
)

// Hub represents a interface for the backend hub service.
type Hub struct {
	// HostAddr is the address of the backend hub service.
	HostAddr string

	// AuthToken is the user authentication token.
	AuthToken string

	// conn holds the websocket connection.
	conn *websocket.Conn

	// updates
	updates chan []byte
}

// NewHub creates and start a websocket connection to backend hub.
func NewHub(addr string, authToken string) *Hub {
	hub := &Hub{
		HostAddr:  addr,
		AuthToken: authToken,
		updates:   make(chan []byte),
	}
	return hub
}

// Header generate the hub request header.
func (h *Hub) Header() http.Header {
	header := http.Header{}
	header.Set("user-token", h.AuthToken)
	return header
}

// URL generate the hub request url.
func (h *Hub) URL() url.URL {
	return url.URL{
		Scheme: "ws",
		Host:   h.HostAddr,
		Path:   "/device-sync",
	}
}

// Dial dial the backend hub and establish a websocket connection
// and stores the connection to the hub conn.
func (h *Hub) Dial() error {
	url := h.URL()
	log.Printf("Attempting connection to: %s", url.String())
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), h.Header())
	if err != nil {
		log.Println(err)
		time.Sleep(3 * time.Second)
		return h.Dial()
	}
	h.conn = conn
	return nil
}

// SyncTree listens to incoming traffics from the backend hub and
// call the appropriate handler to mutate the vtree.
func (h *Hub) SyncTree(vt *vtree.VTree) error {
	for {
		msg, err := h.ReadMsg()
		if err != nil {
			sys.Alert(err.Error())
		}
		log.Printf("Sync read message: %s", common.ToStr(msg))
		h.updates <- msg
	}
}

// Stop closes the hub websocket connection to the backend hub.
func (h *Hub) Stop() {
	if !h.isConnected() {
		return
	}
	defer h.conn.Close()
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := h.conn.WriteMessage(websocket.CloseMessage, closeMsg)
	if err != nil {
		sys.Alert(err.Error())
	}
}

// PushMsg send a msg to websocket connection
func (h *Hub) PushMsg(msg []byte) error {
	if !h.isConnected() {
		return ErrHubConnFailed
	}
	return h.conn.WriteMessage(websocket.TextMessage, msg)
}

// ReadMsg read a msg from the websocket connetion
func (h *Hub) ReadMsg() ([]byte, error) {
	if !h.isConnected() {
		return []byte{}, ErrHubConnFailed
	}
	_, msg, err := h.conn.ReadMessage()
	return msg, err
}

func (h *Hub) isConnected() bool {
	if h.conn == nil {
		return false
	}
	return true
}
