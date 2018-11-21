package utils


import (
	"bytes"
	"encoding/gob"
)

func Marshall(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	return buf.Bytes(), err
}

func UnMarshall(data []byte, size int, e interface{}) error {
	var buffer bytes.Buffer
	dec := gob.NewDecoder(&buffer)
	buffer.Write(data[0:size])
	err := dec.Decode(e)
	return err
}