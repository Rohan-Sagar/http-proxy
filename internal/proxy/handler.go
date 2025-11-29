package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net"

	httputils "github.com/rohan-sagar/http-proxy/internal/http"
)

// Forward client requests to backend server
func (p *Proxy) handleRequestForwarding(port string, request *httputils.Request) error {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		log.Printf("Error connecting to server: %v\n", err)
		return err
	}
	log.Printf("Connected to remote address :%s\n", conn.RemoteAddr())
	fmt.Println("Writing request: %s\n", request.ToString())

	n, err := conn.Write(request.ToBytes())
	if err != nil {
		log.Printf("Error writing to server: %v\n", err)
		return err
	}
	fmt.Print("n", n)
	return nil
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
	error := p.handleRequestForwarding(":8081", newRequest) // static port for now
	if error != nil {
		log.Printf("Failed to forward request: %v\n", error)
		return
	}

	response := httputils.NewResponse(200, "Hallelujah")

	conn.Write([]byte(response.ToString()))

}
