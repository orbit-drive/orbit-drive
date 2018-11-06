package fs

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/sys"
	"github.com/wlwanpan/orbit-drive/fs/vtree"
)

type Hub struct {
	HostAddr string
	conn     *websocket.Conn
}

func NewHub(addr string) (*Hub, error) {
	hub := &Hub{HostAddr: addr}
	if err := hub.Connect(); err != nil {
		return &Hub{}, nil
	}
	return hub, nil
}

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

func (h *Hub) Sync(vt *vtree.VTree) {
	for {
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			sys.Alert(err.Error())
		}
		log.Println(common.ToStr(message))
	}
}

func (h *Hub) Stop() {
	h.conn.Close()
}

// Push send a msg to websocket connection
func (h *Hub) Push(msg []byte) error {
	return nil
}
