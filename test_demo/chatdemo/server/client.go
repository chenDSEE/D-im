package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	server *Server

	// client info
	name string
	status string	// UP\DOWN

	// manage
	mtx sync.Mutex
	conn *websocket.Conn
	pending int	// lost PONG counter
}

func NewClient(server *Server, conn *websocket.Conn) *Client {
	return &Client{
		server: server,
		status: "DOWN",
		conn: conn,
		pending: 0,
	}
}

func (client *Client) SetName(name string) {
	client.mtx.Lock()
	defer client.mtx.Unlock()

	client.name = name
}

func (client *Client) Status() string {
	client.mtx.Lock()
	defer client.mtx.Unlock()

	return fmt.Sprintf("%s(%d)", client.status, client.pending)
}

func (client *Client) ReadMessage() (int, []byte, error) {
	return client.conn.ReadMessage()
}

func (client *Client) WriteMessage(messageType int, data []byte) error {
	return client.conn.WriteMessage(messageType, data)
}

func (client *Client) Addr() string {
	client.mtx.Lock()
	defer client.mtx.Unlock()

	if client.conn == nil {
		return "NULL"
	}

	return client.conn.RemoteAddr().String()
}

func (client *Client) StatusUp() {
	client.mtx.Lock()
	defer client.mtx.Unlock()

	client.status = "UP"
	client.pending = 0
}

func (client *Client) StatusDown() {
	client.mtx.Lock()
	defer client.mtx.Unlock()

	client.status = "DOWN"
}

func (client *Client) Close() {
	client.mtx.Lock()
	defer client.mtx.Unlock()

	if client.conn != nil {
		_ = client.conn.Close()
	}
}