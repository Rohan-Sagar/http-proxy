package http

import (
	"bufio"
	"errors"
	"io"
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

func cleanString(s string) string {
	return strings.ReplaceAll(s, "\r\n", "")
}
