package main

import (
	"fmt"
	"net"
	"sort"
)

const (
	MAX_PORTS = 65535
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, MAX_PORTS)
	results := make(chan int)
	var open_ports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= MAX_PORTS; i++ {
			ports <- i
		}
	}()

	for i := 0; i < MAX_PORTS; i++ {
		port := <-results
		if port != 0 {
			open_ports = append(open_ports, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(open_ports)
	for _, port := range open_ports {
		fmt.Printf("%d open\n", port)
	}
}
