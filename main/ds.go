package main

import (
	"../dirserver"
	"log"
	"os"
)

func main() {

	Ip := "localhost"
	PortForTN := "8001"
	PortForTC := "8002"
	PortForHB := "8003"

	if len(os.Args) == 5 {
		Ip = os.Args[1]
		PortForTN = os.Args[2]
		PortForTC = os.Args[3]
		PortForHB = os.Args[4]
	} else if len(os.Args) != 1 {
		log.Fatal("usage: go run ds.go [Ip] [PortForTN] [PortForTC] [PortForHB]")
	}

	dirserver.StartDS(Ip, PortForTN, PortForTC, PortForHB)
}
