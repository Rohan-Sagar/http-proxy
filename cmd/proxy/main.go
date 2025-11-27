package main

import (
	"fmt"
	"log"
	"net"
)

var PORT = "8080"

func handleConn(conn net.Conn) {
	defer func() {
		log.Printf("Connection closed %s\n", conn.RemoteAddr())
		conn.Close()
	}()
}

func main() {
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started on port:", PORT)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Printf("Accepted connection from: %s\n", conn.RemoteAddr())
		go handleConn(conn)
	}
}
