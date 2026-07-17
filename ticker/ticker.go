package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("tick!")
			case <-done:
				fmt.Println("done")
				return
			}
		}
	}()

	time.Sleep(5 * time.Second) // let it tick for 5 seconds
	done <- true                // we'll use this next — hold that thought
	fmt.Println("stopped")
}
