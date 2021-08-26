package main

import (
	"fmt"
	"log"
	"time"
)

const DELAY_TIME = 1	// second

func main() {
	count := 0
	// ticker 是周而复始的；timer 则是一次性的
	ticker := time.NewTicker(time.Second * DELAY_TIME)

	// start tick
	go func() {
		for {
			select {
			case <-ticker.C:
				count++;
				fmt.Printf("tick, %d\n", count)
				if count == 10 {
					log.Fatal("time up, shutdown !!!")
				}
			}
		}
	}()

	time.Sleep(time.Second * 20)
}
