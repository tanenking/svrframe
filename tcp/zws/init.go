package zws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader websocket.Upgrader
	cID      uint32
)

func checkOrigin(r *http.Request) bool {
	return true
}

func init() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     checkOrigin,
	}
	cID = 0
}
