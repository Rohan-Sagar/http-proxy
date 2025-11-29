package proxy

import (
	"fmt"
	"log"
	"net"

	"github.com/rohan-sagar/http-proxy/internal/config"
)

type Proxy struct {
	config   *config.Config
	listener net.Listener
}

func New(cfg *config.Config) *Proxy {
	return &Proxy{
		config: cfg,
	}
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", ":"+p.config.Port)
	if err != nil {
		return err
	}
	p.listener = listener

	fmt.Printf("Server started on port: %s\n", p.config.Port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Printf("Accepted connection from: %s\n", conn.RemoteAddr())
		go p.handleConnection(conn)
	}
}
