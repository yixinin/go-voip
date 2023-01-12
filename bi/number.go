package bi

import (
	"bytes"
	"encoding/binary"
)

func BytesToInt[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](data []byte) T {
	var i T
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, &i); err != nil {
		panic(err)
	}
	return i
}
