package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/DistributedClocks/GoVector/govec"
)

const MSG_SIZE = 4                             // 4 bytes in size

func TCPRead(from *net.TCPConn, vecLogger *govec.GoLog, vecMsg string) ([]byte, error) {
	sizeBuf := make([]byte, MSG_SIZE)
	_, err := from.Read(sizeBuf)

	if err != nil {
		return nil, err
	}

	mlen := binary.LittleEndian.Uint32(sizeBuf)

	if err != nil {
		return nil, err
	}

	actlen := int(mlen)
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

	b := make([]byte, MSG_SIZE)
	binary.LittleEndian.PutUint32(b, uint32(len(loggedPayload)))

	_, _  = to.Write(b)

	return to.Write(loggedPayload)
}
