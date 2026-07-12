package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

type ScanResults struct {
	Port int
	Open bool
}

func workers(host string, ports <-chan int, mu *sync.Mutex, results *[]ScanResults, wg *sync.WaitGroup) {
	for port := range ports {
		func() {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", host, port)
			conn, err := net.DialTimeout("tcp", address, time.Second)
			result := ScanResults{Port: port}
			if err == nil {
				conn.Close()
				result.Open = true
			}
			mu.Lock()
			*results = append(*results, result)
			mu.Unlock()
		}()
	}
}

func main() {
	host := "scanme.nmap.org"
	start := 20
	end := 500
	numWorkers := 500

	var wg sync.WaitGroup
	var mu sync.Mutex
	var collected []ScanResults
	ports := make(chan int, end-start+1)

	for range numWorkers {
		go workers(host, ports, &mu, &collected, &wg)
	}

	for port := start; port <= end; port++ {
		wg.Add(1)
		ports <- port
	}

	close(ports)

	wg.Wait()

	sort.Slice(collected, func(i, j int) bool {
		return collected[i].Port < collected[j].Port
	})

	for _, r := range collected {
		status := "close"
		if r.Open {
			status = "open"
		}
		fmt.Printf("%d: %s\n", r.Port, status)
	}

}
