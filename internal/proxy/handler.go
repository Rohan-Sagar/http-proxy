package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net"

	httputils "github.com/rohan-sagar/http-proxy/internal/http"
)

// Handle response from backend server
func (p *Proxy) handleServerResponse(conn net.Conn) {
	fmt.Print("handleServerResponse\n")
	reader := bufio.NewReader(conn)
	response, err := httputils.ParseResponse(reader)
	if err != nil {
		log.Printf("Failed to parse response body: %v\n", err)
		return
	}
	fmt.Println("response", response.ToString())
}

// Forward client requests to backend server
func (p *Proxy) handleRequestForwarding(port string, request *httputils.Request) (conn net.Conn, error error) {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		log.Printf("Error connecting to server: %v\n", err)
		return conn, err
	}
	log.Printf("Connected to remote address :%s\n", conn.RemoteAddr())
	fmt.Println("Writing request: %s\n", request.ToString())

	_, err = conn.Write(request.ToBytes()) // _ is the number of bytes written to the server - equal to request.ToBytes
	if err != nil {
		log.Printf("Error writing to server: %v\n", err)
		return conn, err
	}
	return conn, nil
}

/*
Processes a single client connection:
1) Parses incoming HTTP request
2) Sends a response back to the client
*/
func (p *Proxy) handleConnection(conn net.Conn) {
	incomingAddress := conn.RemoteAddr()
	defer func() {
		log.Printf("Connection closed %s\n", incomingAddress)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	request, err := httputils.ParseRequest(reader)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		return
	}

	log.Printf("Request: %+v\n", request)

	newRequest := httputils.NewRequest(request, incomingAddress)
	serverConn, error := p.handleRequestForwarding(":8081", newRequest) // static port for now
	if error != nil {
		log.Printf("Failed to forward request: %v\n", error)
		return
	}

	p.handleServerResponse(serverConn)
}
