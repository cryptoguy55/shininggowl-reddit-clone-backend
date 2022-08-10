package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	EventTypeMsg    = "event-msg"    // 用户发言
	EventTypeSystem = "event-system" // 系统信息推送 如房间人数
	EventTypeJoin   = "event-join"   // 用户加入
	EventTypeTyping = "event-typing" // 用户正在输入
	EventTypeLeave  = "event-leave"  // 用户离开
	EventTypeImage  = "event-image"  // todo 消息图片
)

type response struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Type      string `json:"type,omitempty"`
}

var Manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true

			jsonMessage, _ := json.Marshal(&Message{Content: "A new socket has connected.", Type: EventTypeJoin, Sender: conn.id})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected.", Type: EventTypeLeave, Sender: conn.id})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {

				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		fmt.Println(conn.id)
		if conn != ignore {
			conn.send <- message
		}
	}
}
func (c *Client) read() {
	defer func() {
		Manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		res := response{}
		json.Unmarshal([]byte(string(message)), &res)
		fmt.Println(string(message))
		if err != nil {
			Manager.unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Recipient: res.Recipient, Content: res.Content, Type: EventTypeMsg, Timestamp: (time.Now()).Unix() * 1000})
		Manager.broadcast <- jsonMessage
	}
}
func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
func WsPage(res http.ResponseWriter, req *http.Request, id string) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &Client{id: id, socket: conn, send: make(chan []byte)}

	Manager.register <- client
	fmt.Println(client)
	go client.read()
	go client.write()
}
