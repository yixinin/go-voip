package av

type Alive interface {
	Alive() bool
}

type Closer interface {
	Close()
}

type AvReader interface {
	Closer
	Alive
	Read(*Packet)
}

type AvWriter interface {
	Closer
	Alive
	Write(*Packet)
}
