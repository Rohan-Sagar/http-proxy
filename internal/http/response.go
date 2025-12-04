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
	RawHeader  string
}

// For debugging - read the response as a string
func (r *Response) ToString() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, r.StatusText))

	for k, v := range r.Headers {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	sb.WriteString("\r\n")
	sb.WriteString("(body streamed, not buffered)\n")

	return sb.String()
}

// Parse HTTP/1.1 response from raw TCP connection
func ParseResponseHeaders(reader *bufio.Reader) (*Response, error) {
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

	statusCode, err := strconv.Atoi(cleanString(requestLine[1]))
	if err != nil {
		log.Printf("Failed to convert statusCode string to int: %v\n", err)
		return nil, err
	}
	headers := parseHeaders(lines)

	return &Response{
		StatusCode: statusCode,
		StatusText: cleanString(requestLine[2]),
		Headers:    headers,
		RawHeader:  buildRawHeaders(lines),
	}, nil
}
