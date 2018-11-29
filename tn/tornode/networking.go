package tornode

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

func tcpRead(from *net.TCPConn) ([]byte, error) {
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

func tcpWrite(to *net.TCPConn, payload []byte) (int, error) {
	return to.Write(payload)
}

// GetNewUnusedPort generate a random local ip port
func getNewUnusedPort() string {
	// TODO: you should check whether this random port is really unused
	// assuming IPs are all localhost
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	// generate a random port between [min, max]
	min := 3000
	max := 60000
	port := random.Intn(max-min) + min
	return strconv.Itoa(port)
}
