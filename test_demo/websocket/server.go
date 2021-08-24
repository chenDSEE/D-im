package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)


func echoHandle(resp http.ResponseWriter, req *http.Request) {
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

	defer web.Close()

	// websocket service as a long connect
	for {
		// receive data
		messageType, data, err := web.ReadMessage()
		if err != nil {
			log.Println("read message from websocket fail !!!", err)
			return
		}
		// messageType is websocket.TextMessage(1)
		log.Printf("[%s-->%s][Type:%d], data:[%s]\n",
			       web.RemoteAddr().String(), web.LocalAddr().String(),
			       messageType, data)

		// do echo
		err = web.WriteMessage(messageType, data)
		if err != nil {
			log.Println("write message to websocket fail !!!", err)
			return
		}
	}
}

func main() {
	log.Println("==== server start ====")

	// depoly HTTP server
	http.HandleFunc("/echo", echoHandle)
	if err := http.ListenAndServe("localhost:8003", nil); err != nil {
		log.Fatal(err)
	}

	log.Println("==== server stop ====")
}
