package rw

const (
	ProtocolTCP  = "tcp"
	ProtocolUDP  = "udp"
	ProtocolWS   = "ws"
	ProtocolHttp = "http"
)

type Closer interface {
	Close() error
}

type ReaderCloser interface {
	Closer
	Read([]byte) (int, error)
}

type WriterCloser interface {
	Closer
	Write([]byte) (int, error)
}

type ReaderWriterCloser interface {
	Closer
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Name() string
}
