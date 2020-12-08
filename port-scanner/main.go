package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

func worker(host string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
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
	var portsRange, host string
	flag.StringVar(&portsRange, "ports", "1-1024", "port range to be scanned")
	flag.StringVar(&host, "host", "scanme.nmap.org", "target host")
	flag.Parse()

	limits := strings.Split(portsRange, "-")
	start, _ := strconv.Atoi(limits[0])
	end, _ := strconv.Atoi(limits[1])

	log.Printf("Scannning ports from %d to %d in %s...\n", start, end, host)

	ports := make(chan int, 100)
	results := make(chan int)
	var openPorts []int

	for i := 0; i < cap(ports); i++ {
		go worker(host, ports, results)
	}

	go func() {
		for i := start; i <= end; i++ {
			ports <- i
		}
	}()

	for i := start; i <= end; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}
	close(ports)
	close(results)

	if len(openPorts) == 0 {
		log.Printf("No open ports found between %d and %d...\n", start, end)

	} else {
		sort.Ints(openPorts)
		for _, p := range openPorts {
			fmt.Println("Open port", p)
		}
	}
}
