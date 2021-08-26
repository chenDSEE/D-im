package main

import (
	"log"
	"chatdemo/server"
	"time"
)

// FEATURE:
// 1. same server just like a chat-room, all the word you said, will be broadcast
// 2. same user overwrite the oldder one
// 3. healthy check(if timeout, server clear this subscriber)
//    3.1 you can buffer all the data for 1 min, after 1 min, clear them.
//        not now, it is not good to buffer message in server.
// 4. client auto reconnect(server do nothing)
// 5. broadcast user log-out information to others

// USAGE:
// ws://192.168.8.129:8003/?user=aimer
// ws://192.168.8.129:8003/?user=chelly

func main() {

	// TODO: use a struct 'config' to to new a server
	ser := server.NewServer("192.168.8.129:8003")

	go func() {
		for {
			// show server status every 10 seconds
			time.Sleep(10 * time.Second)
			ser.ShowServerStatus()
		}
	}()

	if err := ser.Start(); err != nil {
		log.Printf("server error with [%v]\n", err)
	}

	log.Println("==== server stop ====")
}
