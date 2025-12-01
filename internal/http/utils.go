package http

import (
	"bufio"
	"strings"
)

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

// remove \r\n from strings
func cleanString(s string) string {
	return strings.ReplaceAll(s, "\r\n", "")
}
