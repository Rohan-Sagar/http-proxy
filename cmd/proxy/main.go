package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var PORT = "8080"

type RequestMethod string

const (
	Get     RequestMethod = "GET"
	Put     RequestMethod = "PUT"
	Post    RequestMethod = "POST"
	Patch   RequestMethod = "PATCH"
	Delete  RequestMethod = "DELETE"
	Head    RequestMethod = "HEAD"
	Options RequestMethod = "OPTIONS"
)

type Request struct {
	Method          RequestMethod
	Path            string
	ProtocolVersion string
	Headers         map[string]string // using general map instead of Headers struct since we should be able to forward unknown headers
}

func cleanString(s string) string {
	substring := "\r\n"
	// if the string contains the end of line delimiter - remove
	s = strings.ReplaceAll(s, substring, "")
	return s
}

func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range lines[1:] {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0])) // lower case for ease
		val := strings.TrimSpace(parts[1])

		headers[key] = val
	}
	return headers
}

func handleConn(conn net.Conn) {
	defer func() {
		log.Printf("Connection closed %s\n", conn.RemoteAddr())
		conn.Close()
	}()
	// parse HTTP/1.1 reqs from raw TCP connections
	reader := bufio.NewReader(conn)
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading request string: %v\n", err)
			break
		}
		// end of headers
		if line == "\r\n" {
			break
		}
		fmt.Print(line)
		lines = append(lines, line)
	}

	headers := parseHeaders(lines)
	requestLine := strings.SplitN(lines[0], " ", 3)
	request := Request{
		Method:          RequestMethod(requestLine[0]),
		Path:            requestLine[1],
		ProtocolVersion: cleanString(requestLine[2]),
		Headers:         headers,
	}

	fmt.Println("\n")
	fmt.Printf("%+v\n", request)
	fmt.Println("\n")

	fmt.Print(reader.)
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
