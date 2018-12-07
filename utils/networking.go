package utils

import (
	"net"

	"github.com/DistributedClocks/GoVector/govec"
)

func TCPRead(from *net.TCPConn, vecLogger *govec.GoLog, vecMsg string) ([]byte, error) {
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
	var results []byte
	vecLogger.UnpackReceive(vecMsg, bytes, &results, govec.GetDefaultLogOptions())
	return results, nil
}

func TCPWrite(to *net.TCPConn, payload []byte, vecLogger *govec.GoLog, vecMsg string) (int, error) {
	loggedPayload := vecLogger.PrepareSend(vecMsg, payload, govec.GetDefaultLogOptions())
	return to.Write(loggedPayload)
}
