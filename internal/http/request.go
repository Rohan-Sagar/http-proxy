package http

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

var (
	ErrMalformedRequest = errors.New("malformed request")
)

type Request struct {
	Method          RequestMethod
	Path            string
	ProtocolVersion string
	Headers         map[string]string // using general map instead of Headers struct to forward unknown headers
	Body            []byte
}

/*
Converts a Request struct into raw HTTP request bytes
exactly like the incoming HTTP/1.1 request sent over
a real TCP conn
eg:

	POST /api/users HTTP/1.1 -- request line
	Host: localhost:8080 -- header
	Accept: * / * -- header (without spaces)
	Content-Length: 25 -- header
	Content-Type: application/json -- header
	User-Agent: curl/8.7.1 -- header
	{"name": "Rohan"} -- body
*/
func (r *Request) ToBytes() []byte {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%s %s %s\r\n", r.Method, r.Path, r.ProtocolVersion)
	for k, v := range r.Headers {
		fmt.Fprintf(&b, "%s: %s\r\n", k, v)
	}

	//end of headers
	b.WriteString("\r\n")

	// body if it exists
	if len(r.Body) > 0 {
		b.Write(r.Body)
	}

	return b.Bytes()
}

// Convert Request struct to a string - for debug purposes
func (r *Request) ToString() string {
	resp := fmt.Sprintf("%s %s %s\r\n", r.Method, r.Path, r.ProtocolVersion)

	for key, val := range r.Headers {
		resp += fmt.Sprintf("%s: %s\r\n", key, val)
	}

	resp += "\r\n"
	resp += string(r.Body)

	return resp
}

/*
Construct a new Request to send to the backend with additional headers:
1) X-Forwarded-For - originating IP address of the client
2) X-Forwarded-Proto - http or https
*/
func NewRequest(r *Request, incomingAddress net.Addr) *Request {
	headers := make(map[string]string)

	for k, v := range r.Headers {
		headers[k] = v // this is fine because request headers are case insensitive
	}

	headers["x-forwarded-for"] = incomingAddress.String()
	headers["x-forwarded-proto"] = "http" // static for now

	return &Request{
		Method:          r.Method,
		Path:            r.Path,
		ProtocolVersion: r.ProtocolVersion,
		Headers:         headers,
		Body:            r.Body,
	}
}

// Parse HTTP/1.1 request from raw TCP connection
func ParseRequest(reader *bufio.Reader) (*Request, error) {
	lines, err := readHeaders(reader)
	if err != nil {
		return nil, err
	}

	if len(lines) < 1 {
		return nil, ErrMalformedRequest
	}

	requestLine := strings.SplitN(lines[0], " ", 3)
	if len(requestLine) < 3 {
		return nil, ErrMalformedRequest
	}

	headers := parseHeaders(lines)
	body := parseBody(reader, headers)

	return &Request{
		Method:          RequestMethod(requestLine[0]),
		Path:            requestLine[1],
		ProtocolVersion: cleanString(requestLine[2]),
		Headers:         headers,
		Body:            body,
	}, nil
}

// Read HTTP header lines from the connection until it encounters a blank line - which means reached end of headers according to HTTP/1.1
func readHeaders(reader *bufio.Reader) ([]string, error) {
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// end of headers
		if line == "\r\n" { // \r = carriage return (move to start of line), \n = line feed (move to next line) - comes from typewrites
			break
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// Convert raw header lines into a hashmap
func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range lines[1:] {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0])) // lowercase for ease
		val := strings.TrimSpace(parts[1])
		headers[key] = val
	}
	return headers
}

// Read the length of the Content-Length Header and read exactly those many bytes
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
	io.ReadFull(reader, body)
	return body
}

// remove \r\n from strings
func cleanString(s string) string {
	return strings.ReplaceAll(s, "\r\n", "")
}
