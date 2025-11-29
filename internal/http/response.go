package http

import (
	"bytes"
	"fmt"
)

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

func (r *Response) Byte() []byte {
	var buf bytes.Buffer
	buf.WriteString("\r\n")
	buf.Write(r.Body)
	return buf.Bytes()
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

// creates an http response
func NewResponse(statusCode int, body string) *Response {
	return &Response{
		StatusCode: statusCode,
		StatusText: getStatusText(statusCode),
		Headers: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(body)),
		},
		Body: []byte(body),
	}
}

func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}
