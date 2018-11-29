package main

import (
	"fmt"
	"os"
	"strconv"

	"./tornode"
)

func main() {
	args := os.Args[1:]

	if !(len(args) == 0 || len(args) == 6) {
		fmt.Println("Usage: go run tn/main.go [dsIPPort] [listenIPPort] [fdListenIPPort] [private key path] [public key path] [timeOutMillis]")
		return
	}

	dsIPPort := "127.0.0.1:3000"
	listenIPPort := "127.0.0.1:4001"
	fdListenIPPort := "127.0.0.1:4002"
	privateKeyPath := "./tn/keys/public.pem"
	publicKeyPath := "./tn/keys/private.pem"
	timeOutMillis := 1000

	if len(args) == 6 {
		dsIPPort = args[0]
		listenIPPort = args[1]
		fdListenIPPort = args[2]
		privateKeyPath = args[3]
		publicKeyPath = args[4]
		var err error
		timeOutMillis, err = strconv.Atoi(args[5])
		if err != nil {
			fmt.Println("Invalid timeOutMillis integer")
			return
		}
	}

	fmt.Println("launching tor node")
	tnerr := tornode.InitTorNode(dsIPPort, listenIPPort, fdListenIPPort, privateKeyPath, publicKeyPath, timeOutMillis)
	fmt.Println(tnerr)
}
