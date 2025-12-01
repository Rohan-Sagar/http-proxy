package http

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

/*
Builds an http response string in the exact format required by the HTTP/1.1 protocol
Example:

	"HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/plain\r\n" +
	"Content-Length: 10\r\n" +
	"\r\n" +
	"Hallelujah"
*/
func (r *Response) ToString() string {
	resp := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, r.StatusText)

	for key, val := range r.Headers {
		resp += fmt.Sprintf("%s: %s\r\n", key, val)
	}

	resp += "\r\n"
	resp += string(r.Body)

	return resp
}

// Parse HTTP/1.1 response from raw TCP connection
func ParseResponse(reader *bufio.Reader) (*Response, error) {
	lines, err := readHeaders(reader)
	if err != nil {
		return nil, err
	}

	if len(lines) < 1 {
		return nil, ErrMalformedRequest
	}

	// ProtocolVersion StatusCode StatusText (HTTP/1.1 200 OK)
	requestLine := strings.SplitN(lines[0], " ", 3)
	if len(requestLine) < 3 {
		return nil, ErrMalformedRequest
	}

	headers := parseHeaders(lines)
	body := parseBody(reader, headers)

	statusCode, err := strconv.Atoi(cleanString(requestLine[1]))
	if err != nil {
		log.Printf("Failed to convert statusCode string to int: %v\n", err)
		return nil, err
	}

	return &Response{
		StatusCode: statusCode,
		StatusText: cleanString(requestLine[2]),
		Headers: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(body)),
		},
		Body: []byte(body),
	}, nil
}
