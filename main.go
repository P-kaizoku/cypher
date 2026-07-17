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

type Scanner struct {
	host         string
	results      []ScanResults
	mu           sync.Mutex
	wg           sync.WaitGroup
	scannedPorts int
	scanMu       sync.Mutex
}

func worker(s *Scanner, ports <-chan int) {
	for port := range ports {
		func() {
			defer s.wg.Done()
			address := fmt.Sprintf("%s:%d", s.host, port)
			conn, err := net.DialTimeout("tcp", address, time.Second)
			result := ScanResults{Port: port}
			if err == nil {
				conn.Close()
				result.Open = true
			}
			s.mu.Lock()
			s.results = append(s.results, result)
			s.mu.Unlock()

			s.scanMu.Lock()
			s.scannedPorts++
			s.scanMu.Unlock()
		}()
	}
}

func main() {
	host := flag.String("host", "scanme.nmap.org", "network to scan")
	start := flag.Int("start", 1, "start port")
	end := flag.Int("end", 1024, "end port")
	worker_ := flag.Int("worker", 100, "concurrent workers")
	verbose := flag.Bool("verbose", false, "to show close ports")
	flag.Parse()

	startTime := time.Now()

	if *end < *start {
		fmt.Println("start port must be smaller than ending port")
		return
	}

	s := &Scanner{
		host: *host,
	}

	ports := make(chan int, *end-*start+1)

	for range *worker_ {
		go worker(s, ports)
	}

	for port := *start; port <= *end; port++ {
		s.wg.Add(1)
		ports <- port
	}
	close(ports)

	total := *end - *start + 1
	done := make(chan bool)

	go func() {
		ticker := time.NewTicker(time.Millisecond * 200)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.scanMu.Lock()
				portScanned := s.scannedPorts
				s.scanMu.Unlock()
				fmt.Printf("\r scanned ports: %d/%d", portScanned, total)
			case <-done:
				return
			}
		}
	}()

	s.wg.Wait()
	done <- true
	fmt.Print("\r                                        \r")

	sort.Slice(s.results, func(i, j int) bool {
		return s.results[i].Port < s.results[j].Port
	})

	openPorts := 0
	for _, r := range s.results {
		status := "Close"
		if r.Open {
			status = "Open"
			openPorts++
		}
		if status == "Close" && !*verbose {
			continue
		}
		fmt.Printf("%d: %s\n", r.Port, status)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nscanned %d ports in %v, %d open, %d closed", len(s.results), elapsed, openPorts, len(s.results)-openPorts)
}
