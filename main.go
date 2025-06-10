package main

import (
	"fmt"
	"log"
	"net"
	"sort"
)

const (
	MAX_PORTS uint32 = 99999
)

// https://github.com/thanhphuchuynh/tcp-scanner-go-basic/blob/main/main.go
func worker(ports, results chan int, site string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", site, p)
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
	fmt.Println("Port Scanner - MistoWest")
	fmt.Println("-----------------------------------------")

	var open_ports []int
	var qnt_ports uint32
	var site string

	fmt.Printf("Digite o site: ")
	fmt.Scan(&site)

	fmt.Printf("Quantidade de portas (Máximo - %d): ", MAX_PORTS)
	fmt.Scan(&qnt_ports)

	if qnt_ports > MAX_PORTS {
		log.Fatal("Quantidade de portas maior que o limite")
	}

	if qnt_ports <= 0 {
		log.Fatal("Quantidade de portas menor ou igual a 0")
	}

	ports := make(chan int, qnt_ports)
	results := make(chan int)

	fmt.Printf("Iniciando escaneamento...\n")
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, site)
	}

	go func() {
		for i := 1; i <= int(MAX_PORTS); i++ {
			ports <- i
		}
	}()

	for i := 0; i < int(MAX_PORTS); i++ {
		port := <-results
		if port != 0 {
			open_ports = append(open_ports, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(open_ports)
	fmt.Println("Site: ", site)
	fmt.Printf("Quantidade de portas abertas: %d\n", len(open_ports))
	fmt.Println("-----------------------------------------")
	for _, port := range open_ports {
		fmt.Printf("%d open\n", port)
	}
}
