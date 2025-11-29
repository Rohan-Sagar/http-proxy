package main

import (
	"log"

	"github.com/rohan-sagar/http-proxy/internal/config"
	"github.com/rohan-sagar/http-proxy/internal/proxy"
)

func main() {
	cfg := config.Load()

	p := proxy.New(cfg)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
