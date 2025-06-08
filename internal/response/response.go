package response

import (
	"fmt"
	"io"
	"strconv"
	"streams/internal/headers"
)

type StatusCode int

const (
	Ok          StatusCode = 200
	BadRequest             = 400
	ServerError            = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	mapping := map[StatusCode]string{
		200: "HTTP/1.1 200 OK",
		400: "HTTP/1.1 400 Bad Request",
		500: "HTTP/1.1 500 Internal Server Error",
	}

	msg, present := mapping[statusCode]
	if !present {
		s := fmt.Sprintf("HTTP/1.1 %d Unknown Status\r\n", statusCode)

		_, err := w.Write([]byte(s))
		return err
	}

	_, err := w.Write([]byte(msg + "\r\n"))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	contentLenStr := strconv.Itoa(contentLen)
	return headers.Headers{
		"Content-Length": contentLenStr,
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		h := fmt.Sprintf("%s: %s\r\n", k, v)
		if _, err := w.Write([]byte(h)); err != nil {
			return err
		}
	}
	// Important: terminate header section
	_, err := w.Write([]byte("\r\n"))
	return err
}
