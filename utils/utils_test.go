package utils

import (
	"fmt"
	"testing"
)

func Test_Bytes(t *testing.T) {
	var b = []byte{1, 2, 11, 12, 13, 14}
	var a = []byte{3, 4, 5, 6, 7, 8, 9, 10}
	// b := Int64ToBytes(1024024)
	var c = make([]byte, len(b)+8)
	copy(c[8+2:], b[2:])
	copy(c[2:8+2], a)
	copy(c[:2], b[:2])
	fmt.Println(c)
}

func Test_Atom(t *testing.T) {
	for i := 0; i < 10; i++ {
		rid := GetRoomID()
		fmt.Println(rid)
	}
}
