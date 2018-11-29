package utils

import (
	"bufio"
	"fmt"
	"net"
)

func WriteToConnection(conn *net.TCPConn, json string) (int, error) {
	n, werr := fmt.Fprintf(conn, json+"\n")
	return n, werr
}

func ReadFromConnection(conn *net.TCPConn) (string, error) {
	json, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return json, nil
}

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
