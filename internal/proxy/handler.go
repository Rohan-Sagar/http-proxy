package proxy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	httputils "github.com/rohan-sagar/http-proxy/internal/http"
)

// Handle response from backend server
func (p *Proxy) handleServerResponse(serverConn net.Conn, clientConn net.Conn) {
	defer func() {
		log.Printf("Server connection closed %s\n", serverConn.RemoteAddr())
		serverConn.Close()
	}()
	serverReader := bufio.NewReader(serverConn) // read response from the servver
	response, err := httputils.ParseResponseHeaders(serverReader)
	if err != nil {
		log.Printf("Failed to parse response body: %v\n", err)
		return
	}

	fmt.Println("==== Server Response Headers ====")
	fmt.Print(response.ToString())

	if _, err := clientConn.Write([]byte(response.RawHeader)); err != nil {
		log.Printf("Failed to write headers to client: %v\n", err)
		return
	}
	if lenStr, ok := response.Headers["content-length"]; ok {
		n, err := strconv.Atoi(lenStr)
		if err != nil {
			log.Printf("Invalid content-length: %v\n", err)
		} else {
			if _, err := io.CopyN(clientConn, serverReader, int64(n)); err != nil {
				log.Printf("Failed to stream body to client: %v\n", err)
			}
			return
		}
	}
	// if chunked or unknown length
	io.Copy(clientConn, serverReader)
}

// Forward client requests to backend server
func (p *Proxy) handleRequestForwarding(addr string, request *httputils.Request) (conn net.Conn, error error) {
	conn, err := net.Dial("tcp", addr) // addr is host:port
	if err != nil {
		log.Printf("Error connecting to server: %v\n", err)
		return conn, err
	}
	log.Printf("Connected to remote address :%s\n", conn.RemoteAddr())
	fmt.Printf("Writing request: %s\n", request.ToString())

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
func (p *Proxy) handleConnection(clientConn net.Conn) {
	incomingAddress := clientConn.RemoteAddr()
	defer func() {
		log.Printf("Client connection closed %s\n", incomingAddress)
		clientConn.Close()
	}()

	reader := bufio.NewReader(clientConn)

	request, err := httputils.ParseRequest(reader)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		return
	}

	log.Printf("Request: %+v\n", request)

	newRequest := httputils.NewRequest(request, incomingAddress)
	serverConn, error := p.handleRequestForwarding(p.config.BackendURL, newRequest) // static port for now
	if error != nil {
		log.Printf("Failed to forward request: %v\n", error)
		return
	}

	p.handleServerResponse(serverConn, clientConn)
}
