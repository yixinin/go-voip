package rw

import (
	"io"
	"net/http"
)

type HttpReaderWriter struct {
	w      http.ResponseWriter
	r      io.ReadCloser
	closed bool
}

func NewHttpReaderWriter(w http.ResponseWriter, r io.ReadCloser) *HttpReaderWriter {
	return &HttpReaderWriter{
		w: w,
		r: r,
	}
}

func (c *HttpReaderWriter) Read(buf []byte) (n int, err error) {

	return c.r.Read(buf)
}

func (c *HttpReaderWriter) Write(buf []byte) (n int, err error) {

	return c.w.Write(buf)
}

func (c *HttpReaderWriter) Close() error {
	if c.closed {
		c.closed = true
		return c.r.Close()
	}
	return nil
}

func (*HttpReaderWriter) Name() string {
	return ProtocolHttp
}
