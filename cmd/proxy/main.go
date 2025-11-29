package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
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
	Body            []byte
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

func parseBody(reader *bufio.Reader, headers map[string]string) []byte {
	bodySizeStr := headers["content-length"]
	if bodySizeStr == "" {
		return nil
	}
	bodySize, err := strconv.Atoi(bodySizeStr)
	if err != nil {
		return nil
	}
	body := make([]byte, bodySize)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
	}
	return body
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
	body := parseBody(reader, headers)
	requestLine := strings.SplitN(lines[0], " ", 3)
	if len(lines) < 3 {
		log.Printf("Malformed request line: %s", lines[0])
		return
	}
	request := Request{
		Method:          RequestMethod(requestLine[0]),
		Path:            requestLine[1],
		ProtocolVersion: cleanString(requestLine[2]),
		Headers:         headers,
		Body:            body,
	}
	fmt.Println("\n")
	fmt.Printf("%+v\n", request)
	fmt.Println("\n")

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 10\r\n" +
		"\r\n" +
		"Hallelujah"

	conn.Write([]byte(response))
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
