package rw

import (
	"io"
	"net/http"
)

type HttpWriter struct {
	io.ReadCloser
	w      http.ResponseWriter
	closed bool
}

func NewHttpReaderWriter(w http.ResponseWriter, r io.ReadCloser) *HttpWriter {
	return &HttpWriter{
		ReadCloser: r,
		w:          w,
	}
}

func (c *HttpWriter) Write(buf []byte) (n int, err error) {
	n, err = c.w.Write(buf)
	if err != nil {
		return n, err
	}
	c.w.(http.Flusher).Flush()
	return n, err
}

func (c *HttpWriter) Close() error {
	if !c.closed {
		c.closed = true
		return c.ReadCloser.Close()
	}
	return nil
}

func (*HttpWriter) Name() string {
	return ProtocolHttp
}
