package converter

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

/*
	BytesToString is to convert slice to string with zero-copy.
*/
func BytesToString(b []byte) string {
	sliceheader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceheader.Data,
		Len:  sliceheader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

/*
	StringToBytes is to convert string to slice with zero-copy.
*/
func StringToBytes(s string) []byte {
	stringheader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh := reflect.SliceHeader{
		Data: stringheader.Data,
		Len:  stringheader.Len,
		Cap:  stringheader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&sh))
}

type Number interface {
	bool | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

func Int2Bytes[V Number](value V) []byte {
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, value)
	return buffer.Bytes()
}

func Bytes2Int[V Number](b []byte, order binary.ByteOrder) (V, error) {
	buffer := bytes.NewBuffer(b)
	var value V
	err := binary.Read(buffer, order, &value)
	if nil != err {
		return value, err
	}
	return value, nil
}
