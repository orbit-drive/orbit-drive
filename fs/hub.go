package fs

import (
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/sys"
	"github.com/wlwanpan/orbit-drive/fs/vtree"
)

type Hub struct {
	conn *websocket.Conn

	State chan vtree.State
}

func NewHub() *Hub {
	return &Hub{
		State: make(chan vtree.State),
	}
}

func (h *Hub) connect() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:4000",
		Path:   "/account/subscribe",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		sys.Notify(err.Error())
	}
	h.conn = conn
}

func (h *Hub) Start() {
	h.connect()
	for {
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			sys.Alert(err.Error())
		}
		h.State <- vtree.State{
			Path: common.ToStr(message),
			Op:   "Update",
		}
	}
}

func (h *Hub) Stop() {
	h.conn.Close()
	close(h.State)
}

// Push send a msg to websocket connection
func (h *Hub) Push(msg []byte) error {
	return nil
}
