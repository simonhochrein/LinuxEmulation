package helpers

import (
	"bytes"
	"encoding/binary"
)

func EncodeUint64(n uint64) []byte {
	data := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(data, n)
	return data
}

func EncodeString(str string) []byte {
	return []byte(str)
}

func PackStruct(i interface{}) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, i)
	return buf.Bytes()
}
