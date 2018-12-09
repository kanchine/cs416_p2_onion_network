package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/DistributedClocks/GoVector/govec"
)

const MSG_SIZE = 4                             // 4 bytes in size
const BufSize = 1024

type Size struct {
	Size int
}

func TCPRead(from *net.TCPConn, vecLogger *govec.GoLog, vecMsg string) ([]byte, error) {
	sizeBuf := make([]byte, BufSize)
	n, err := from.Read(sizeBuf)

	if err != nil {
		return nil, err
	}

	var s Size

	err = json.Unmarshal(sizeBuf[:n], &s)

	if err != nil {
		return nil, err
	}
	//mlen := binary.LittleEndian.Uint32(sizeBuf)

	//actlen := int(mlen)
	actlen := s.Size
	fmt.Println("Message size:", actlen)

	//sizeMsg, err := conn.Read(msgBuff)
	bytes := make([]byte, 0)
	chunkCap := 1024
	chunk := make([]byte, chunkCap)

	sizeMsg := 0

	for {
		size, rerr := from.Read(chunk)
		if rerr != nil {
			if rerr != io.EOF {
				return nil, rerr
			}

			break
		}
		bytes = append(bytes, chunk[:size]...)
		sizeMsg += size
		fmt.Printf("====Read chunk: %d\n", size)

		if sizeMsg == actlen {
			break
		}
	}

	if sizeMsg != actlen {
		return nil, errors.New("msg size wrong")
	}

	var results []byte
	vecLogger.UnpackReceive(vecMsg, bytes, &results, govec.GetDefaultLogOptions())
	return results, nil
}

func TCPWrite(to *net.TCPConn, payload []byte, vecLogger *govec.GoLog, vecMsg string) (int, error) {
	loggedPayload := vecLogger.PrepareSend(vecMsg, payload, govec.GetDefaultLogOptions())

	//b := make([]byte, MSG_SIZE)
	//binary.LittleEndian.PutUint32(b, uint32(len(loggedPayload)))

	size := Size{Size: len(loggedPayload)}

	b, err := json.Marshal(size)

	if err != nil {
		return 0, err
	}

	_, _  = to.Write(b)

	return to.Write(loggedPayload)
}
