package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	host := "scanme.nmap.org"
	start := 20
	end := 500

	var wg sync.WaitGroup

	for port := start; port <= end; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", host, port)

			conn, err := net.DialTimeout("tcp", address, time.Second)
			if err != nil {
				fmt.Printf("%d not open\n", port)
				return
			}
			conn.Close()

			fmt.Printf("%d port open\n", port)

		}(port)
	}
	wg.Wait()
	fmt.Println("Scan completed")
}
