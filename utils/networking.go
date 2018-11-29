package utils

import (
	"net"
)

func TCPRead(from *net.TCPConn) ([]byte, error) {
	bytes := make([]byte, 0)
	chunkCap := 1024
	chunk := make([]byte, chunkCap)

	for {
		size, rerr := from.Read(chunk)
		if rerr != nil {
			return nil, rerr
		}
		bytes = append(bytes, chunk[:size]...)
		if size < chunkCap {
			break
		}
	}
	return bytes, nil
}

func TCPWrite(to *net.TCPConn, payload []byte) (int, error) {
	return to.Write(payload)
}
