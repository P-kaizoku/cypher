package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	host := "scanme.nmap.org"
	start := 20
	end := 100

	for port := start; port <= end; port++ {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, time.Second)
		if err != nil {
			fmt.Println("not open: ", port)
			continue
		}

		conn.Close()
		fmt.Println("open: ", port)
	}
}
