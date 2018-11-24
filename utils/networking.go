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
