package utils

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

func SetInterfacePtr(p interface{}, v interface{}) {
	rv := reflect.ValueOf(p)
	rt := reflect.TypeOf(p)

	if rt.Kind() != reflect.Ptr {
		panic("only support pointer type")
	}

	rv.Elem().Set(reflect.ValueOf(v))
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()[:4]
}

func ConcatBytes(a []byte, b []byte) []byte {
	al := len(a)
	target := make([]byte, al + len(b))
	for i, bt := range a {
		target[i] = bt
	}
	for i, bt := range b {
		target[i + al] = bt
	}

	return target
}