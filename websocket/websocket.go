package ws

import (
	"net/http"
	"database/sql"
	"fmt"

	"github.com/gorilla/websocket"

	message "../message"
)

type Hub struct {
	Connection map[Client]*websocket.Conn
}

type Client struct {
	Conn *websocket.Conn
}


func Cast(db *sql.DB, w http.ResponseWriter, r *http.Request, hub *Hub, c Client) {
	go c.readMsg(db, hub)
}

func unRgister(hub *Hub, c Client) {
	// register from hub
	delete(hub.Connection, c)
	Connection(hub)
}

func Connection(hub *Hub) {
	for _, conn := range hub.Connection {
		var pmMsg []byte
		if len(hub.Connection) == 2 {
			pmMsg = []byte("2")
		} else {
			pmMsg = []byte("1")
		}
		pm, _ := websocket.NewPreparedMessage(websocket.TextMessage, pmMsg)
		conn.WritePreparedMessage(pm)
	}
}

func (c Client) readMsg(db *sql.DB, hub *Hub) {
	defer unRgister(hub, c)
	for {
		msg := message.PostMessage{}
		err := c.Conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println(err)
				break
			}
			res := getRes(db, msg)
			for _, conn := range hub.Connection {
				conn.WriteJSON(*res)
			}
		}
}

func getRes(db *sql.DB, msg message.PostMessage) *message.GetMessage {
	query := "SELECT m.id, m.message, u.id, u.name, u.icon, m.datetime FROM messages as m INNER JOIN users as u ON m.sender_id = u.id WHERE m.id = (SELECT MAX(id) FROM messages) ORDER BY m.id"
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	res := message.New(row)
	return res
}