package tornode

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"../../keyLibrary"
	"../../utils"
)

type TorNode struct {
	PrivateKey    *rsa.PrivateKey
	ListenIPPort  string
	fd            utils.FD
	timeoutMillis int
}

func InitTorNode(dsIPPort string, listenIPPort string, fdListenIPPort string,
	privateKeyPath string, publicKeyPath string, timeoutMillis int) error {

	// Initialize variables
	privateKey, pkerror := keyLibrary.LoadPrivateKey(privateKeyPath)
	if pkerror != nil {
		fmt.Printf("Could not init tor node. Failed to load private key: %s\n", pkerror)
		return pkerror
	}

	// load necessary keys
	publicKey, pubkerror := keyLibrary.LoadPublicKey(publicKeyPath)
	if pubkerror != nil {
		fmt.Printf("Could not init tor node. Failed to load public key: %s\n", pubkerror)
		return pubkerror
	}

	// start failure detector
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	epochNonce := rand.Uint64()
	fd, _, fdliberr := utils.Initialize(epochNonce, 50)
	if fdliberr != nil {
		fmt.Printf("TorNode: failed to start fdlib for error: %s\n", fdliberr)
		return fdliberr
	}
	fdresErr := fd.StartResponding(fdListenIPPort)
	if fdresErr != nil {
		fmt.Printf("TorNode: failed to start responding for error: %s\n", fdresErr)
		return fdliberr
	}

	// join network
	dsstatus, dserror := contactDS(dsIPPort, listenIPPort, fdListenIPPort, publicKey)
	if dserror != nil {
		fmt.Printf("TorNode: Could not contact DS to join tor network for error: %s\n", dserror)
		return dserror
	}
	if !dsstatus {
		return errors.New("TorNode: Network join rejected by DS")
	}

	laddr, laddrErr := net.ResolveTCPAddr("tcp", listenIPPort)
	if laddrErr != nil {
		fmt.Printf("TorNode: Could not resolve listen address for error: %s\n", laddrErr)
		return laddrErr
	}
	listener, lerr := net.ListenTCP("tcp", laddr)
	if lerr != nil {
		fmt.Printf("TorNode: Could not start TCP listening for error: %s\n", lerr)
		return lerr
	}

	go onionHandler(listener, privateKey, timeoutMillis)

	return nil
}
