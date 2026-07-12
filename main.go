package main

import (
	"net"
	"sync"
	"time"
)

type ScanResults struct {
	Port int
	Open bool
}

func workers(host string, ports <-chan int, mu *sync.Mutex, results *[]ScanResults, wg *sync.WaitGroup){
	for port := ports{
		func(){
			defer wg.Done()
			address := fmt.Sprinf("%s:%d", host, port)
			conn, err := net.DialTimeout("tcp", address, time.Second)
			result := ScanResults{Port: port}
			if err == nil {
				conn.Close()
				result.Open = true
			}
		}()
	}
}

func main() {

}
