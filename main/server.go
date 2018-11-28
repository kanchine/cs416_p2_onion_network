package main

import (
	"../DataServer"
	"fmt"
	"log"
	"os"
)
func main() {
	var configFile string

	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else {
		log.Fatal("usage: go run server.go [ConfigFile]")
	}

	server, err := DataServer.Initialize(configFile, "./DataServer/private.pem")

	if err != nil {
		fmt.Println("Server failed to start, exiting...")
		return
	}

	server.StartService()
}