package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)


// always pass by pointer
type Server struct {
	// info
	ipaddr string
	id int

	// session manage
	userPool map[string]*Client // username --> Client
	userMtx	sync.Mutex       // mutex for userPool
	userCount uint64            // should be atomic

	// shutdown control
	downOnce	sync.Once	// for server shutdown

	// clear up
	failing chan *Client
}



func NewServer(ipaddr string) *Server {
	return &Server{
		ipaddr: ipaddr,
		id: 1,
		userPool: make(map[string]*Client, 100),
		userCount: 0,
		failing: make(chan *Client),
	}
}

func (server *Server) Start() error {
	// register the related URL and handler
	http.HandleFunc("/", func (resp http.ResponseWriter, req *http.Request) {
		clientHandler(server, resp, req)
	})

	// http server start to handle request
	log.Printf("==== server%s start ====", server.GetInfo())
	if err := http.ListenAndServe(server.ipaddr, nil); err != nil {
		return err
	}

	return errors.New("server-loop end and exit")
}

func clientHandler(server *Server, resp http.ResponseWriter, req *http.Request) {
	// upgrade to websocket
	// use default parameters for webscoket
	// CheckOrigin field: 防止跨站点伪造请求的攻击
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true	// alway true, just for test
	}}

	web, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Fatal("webscoket upgrade failed !!", err)
		return
	}

	client := NewClient(server, web)

	defer client.Close()

	// update user information to server
	userName := req.URL.Query().Get("user")
	log.Printf("new client[%s][%s] come up", web.RemoteAddr().String(), userName)
	if err = server.addUser(userName, client); err != nil {
		log.Printf("%v", err)
		return
	}
	defer server.delUser(userName)
	client.StatusUp()

	// websocket service as a long connect
	server.serviceUser(userName, client)

}

func (server *Server) serviceUser(name string, client *Client) {
	for {
		// receive data
		messageType, data, err := client.ReadMessage()
		if err != nil {
			log.Println("read message from websocket fail !!!", err)
			return
		}

		// messageType is websocket.TextMessage(1)
		log.Printf("[%s]:[%s]\n", name, data)

		// TODO: do boardcast(need to avoid race condition), in next demo

		// do echo
		err = client.WriteMessage(messageType, data)
		if err != nil {
			log.Println("write message to websocket fail !!!", err)
			return
		}
	}
}



/*	echo service
func serviceUser(name string, web *websocket.Conn, resp http.ResponseWriter, req *http.Request) {
	for {
		// receive data
		messageType, data, err := web.ReadMessage()
		if err != nil {
			log.Println("read message from websocket fail !!!", err)
			return
		}
		// messageType is websocket.TextMessage(1)
		log.Printf("[%s]:[%s]\n", name, data)

		// do echo
		err = web.WriteMessage(messageType, data)
		if err != nil {
			log.Println("write message to websocket fail !!!", err)
			return
		}
	}
}
*/

func (server *Server) GetInfo() string {
	// ignore race condition, no need to be so exactly
	return fmt.Sprintf("[ID:%d][%s]", server.id, server.ipaddr)
}

func (server *Server) addUser(name string, client *Client) error {
	{
		server.userMtx.Lock()
		defer server.userMtx.Unlock()

		_, exist := server.userPool[name]
		if exist {
			tmp := fmt.Sprintf("user[%s] is already existed !!!", name)
			return errors.New(tmp)
		}

		server.userPool[name] = client
		server.userCount += 1
	}

	client.SetName(name)
	return nil
}

func (server *Server) delUser(name string) {
	server.userMtx.Lock()
	defer server.userMtx.Unlock()

	if _, exist := server.userPool[name]; !exist {
		return
	}

	delete(server.userPool, name)
	server.userCount -= 1
}

func (server *Server) ShowServerStatus() {
	server.userMtx.Lock()
	defer server.userMtx.Unlock()

	fmt.Printf("\n\n=============================================\n")
	fmt.Printf("total user: %d\n", server.userCount)
	for name, client := range server.userPool {
		fmt.Printf("[%s][%s][status: %s]\n", client.Addr(), name, client.Status())
	}
	fmt.Printf("=============================================\n\n")
}