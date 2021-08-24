package main

import (
	"log"
	"chatdemo/server"
	"time"
)

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
