package main

import (
	"io"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()

	// create a buffer to store input data
	buf := make([]byte, 512)

	for {
		size, err := conn.Read(buf[0:])
		if err == io.EOF {
			log.Println("Client disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected error")
			break
		}
		log.Printf("Received %d bytes: %s", size, string(buf))

		// send data
		log.Println("Writing data")
		if _, err := conn.Write(buf[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Println("Listening on 0.0.0.0:20080")

	for {
		// wait for connection
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unabel to accept connection")
		}
		go echo(conn)
	}
}
