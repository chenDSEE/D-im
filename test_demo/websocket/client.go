package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

func main() {
	log.Println("==== client start ====")

	// set up target
	u := url.URL{Scheme: "ws", Host: "localhost:8003", Path: "/echo"}

	// start client and dail to server
	web, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("client socket dial fail !!!")
	}
	defer web.Close()

	sig := make(chan struct{})
	go func() {
		sig <- struct{}{}
		for {
			mt, data, err := web.ReadMessage()
			if err != nil {
				log.Fatal("read data fail !!!")
				return
			}

			log.Printf("[%s-->%s][Type:%d], data:[%s]\n",
				web.RemoteAddr().String(), web.LocalAddr().String(),
				mt, data)
		}
	}()

	<-sig	// wait for read function set up

	// send data to server
	for {
		// send data to sever
		err := web.WriteMessage(websocket.TextMessage, []byte(time.Now().String()))
		if err != nil {
			log.Fatal("send data error")
			break
		}

		time.Sleep(time.Second * 1)
	}


	log.Println("==== client stop ====")
}
