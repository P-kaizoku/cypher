package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func worker(host string, ports <-chan int, wg *sync.WaitGroup) {
	for port := range ports {
		func() {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", host, port)
			conn, err := net.DialTimeout("tcp", address, time.Second)
			if err != nil {
				fmt.Printf("%d close\n", port)
				return
			}

			fmt.Printf("%d open\n", port)
			conn.Close()
		}()
	}
}

func main() {
	host := "scanme.nmap.org"
	start := 20
	end := 500
	numWorkers := 600

	ports := make(chan int, end-start+1)
	var wg sync.WaitGroup

	for range numWorkers {
		go worker(host, ports, &wg)
	}

	for port := start; port <= end; port++ {
		wg.Add(1)
		ports <- port
	}

	close(ports)
	wg.Wait()
	fmt.Println("Scan completed")
}
