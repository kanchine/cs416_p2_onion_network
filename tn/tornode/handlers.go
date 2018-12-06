package tornode

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"time"

	"../../utils"
)

// listening for initial onion messages
func onionHandler(listener *net.TCPListener, privateKey *rsa.PrivateKey, timeoutMillis int) {
	for {
		fmt.Printf("TorNode: Waiting for new circuit connection...\n")
		newCircuitConn, aerr := listener.AcceptTCP()
		if aerr != nil {
			fmt.Printf("TorNode: WARNING could not accept an init onion connection: %s\n", aerr)
			continue
		}
		fmt.Printf("TorNode: new circuit connection from %s! kicking off circuit handler...\n", newCircuitConn.RemoteAddr())
		go handleNewCircuitConn(newCircuitConn, privateKey, timeoutMillis)
	}
}

func handleNewCircuitConn(newCircuitConn *net.TCPConn, privateKey *rsa.PrivateKey, timeoutMillis int) {
	rawBytes, rerr := utils.TCPRead(newCircuitConn)

	if rerr != nil {
		fmt.Printf("TorNode: WARNING read from connection error: %s\n", rerr)
		return
	}

	fmt.Printf("TorNode: Received new onion of %d bytes\n", len(rawBytes))

	nextHop, symmKey, payload, peelerr := peelOnion(rawBytes, privateKey)

	if peelerr != nil {
		fmt.Printf("TorNode: WARNING error when peeling onion: %s\n", peelerr)
		return
	}

	fmt.Printf("TorNode: Onion Peel successful\n")

	// todo next: set up conn to next hop, kick off reverse forwarding thread
	laddr, laddrerr := net.ResolveTCPAddr("tcp", ":"+getNewUnusedPort())
	raddr, raddrerr := net.ResolveTCPAddr("tcp", nextHop)
	if laddrerr != nil || raddrerr != nil {
		fmt.Printf("TorNode: WARNING error resolving tcp addr: %s, %s\n", laddrerr, raddrerr)
		return
	}
	nextHopConn, dialerr := net.DialTCP("tcp", laddr, raddr)
	if dialerr != nil {
		fmt.Printf("TorNode: WARNING error dialing next hop: %s\n", dialerr)
		return
	}
	forwardNextHelper(nextHopConn, payload)
	forwardBackHelper(nextHopConn, newCircuitConn, symmKey, timeoutMillis)
}

func forwardNextHelper(to *net.TCPConn, payload []byte) {
	_, err := utils.TCPWrite(to, payload)
	if err != nil {
		fmt.Printf("TorNode: WARNING forward onion to next hop: %s\n", err)
		return
	}
	fmt.Printf("TorNode: Successfully fowarded onion, next hop: %s, payload size: %d\n", to.RemoteAddr(), len(payload))
}

// NOTE: assuming forwarding from Tn+1 to Tn-1
// Note: assuming response data only sent once
func forwardBackHelper(from *net.TCPConn, to *net.TCPConn, symmKey []byte, timeoutMillis int) {
	defer func() {
		from.Close()
		to.Close()
	}()

	derr := from.SetReadDeadline(time.Now().Add(time.Duration(timeoutMillis) * time.Millisecond))
	if derr != nil {
		fmt.Printf("TorNode: WARNING failed to set read deadline: %s\n", derr)
		return
	}
	fmt.Printf("TorNode: Wating response from nextHop: %s\n", from.RemoteAddr())
	payload, rerr := utils.TCPRead(from)

	if dpassederr, ok := rerr.(net.Error); ok && dpassederr.Timeout() {
		fmt.Printf("TorNode: WARNING waiting data from %s timeout. Tearing down forwarding from %s to %s\n", from.RemoteAddr(), from.RemoteAddr(), to.RemoteAddr())
		return
	}

	if rerr == io.EOF {
		fmt.Printf("TorNode: Failed to forward response back: unexpected remote connection from [%s] closed\n", from.RemoteAddr())
		return
	}

	if rerr != nil {
		fmt.Printf("TorNode: WARNING failed to read from connection: %s\n", rerr)
		return
	}

	forwardPayload, oerr := wrapOnion(payload, symmKey)

	if oerr != nil {
		fmt.Printf("TorNode: WARNING could not wrap onion: %s\n", oerr)
		return
	}

	_, werr := utils.TCPWrite(to, forwardPayload)
	if werr != nil {
		fmt.Printf("TorNode: WARNING failed to forward previous hop: %s\n", werr)
		return
	}
	fmt.Printf("TorNode: Successfully fowarded onion from %s BACK to %s, payload size: %d\n", from.RemoteAddr(), to.RemoteAddr(), len(payload))
}
