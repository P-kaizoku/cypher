package main

import (
	"flag"
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
	host := flag.String("host", "scanme.nmap.org", "network to scan")
	start := flag.Int("start", 1, "start port")
	end := flag.Int("end", 1024, "end port")
	worker := flag.Int("workers", 100, "concurrent workers")
	verbose := flag.Bool("verbose", false, "to show close ports")
	flag.Parse()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var collected []ScanResults
	ports := make(chan int, end-start+1)

	for range *worker {
		go workers(*host, ports, &mu, &collected, &wg)
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
		status := "Close"
		if r.Open {
			status = "Open"
		}
		fmt.Printf("%d: %s\n", r.Port, status)
	}

}
