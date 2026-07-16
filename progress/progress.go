package main

import (
	"fmt"
	"time"
)

func main() {
	count := 0
	for {
		fmt.Print("\r") // move cursor to start of line
		fmt.Printf("elapsed: %d seconds", count)
		time.Sleep(time.Second)
		count++
	}
}
