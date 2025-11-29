package proxy

import (
	"bufio"
	"log"
	"net"

	httputil "github.com/rohan-sagar/http-proxy/internal/http"
)

func (p *Proxy) handleConnection(conn net.Conn) {
	defer func() {
		log.Printf("Connection closed %s\n", conn.RemoteAddr())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	request, err := httputil.ParseRequest(reader)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		return
	}

	log.Printf("Request: %+v\n", request)

	response := httputil.NewResponse(200, "Hallelujah")
	conn.Write([]byte(response.String()))
}
