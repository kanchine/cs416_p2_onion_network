package TorClient

import (
	"../utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	configPath := ""

	if len(os.Args) - 1 != 2 {
		fmt.Println("please use: go run client.go client.json keyToFetch")
	}

	configPath = os.Args[1]

	rawConfig, fileerr := ioutil.ReadFile(configPath)
	if fileerr != nil {
		log.Printf("client.go: Invalid json file: %s\n", fileerr)
		os.Exit(1)
	}
	clientConfig := &utils.ClientConfig{}
	jsonErr := json.Unmarshal(rawConfig, clientConfig)
	if jsonErr != nil {
		log.Printf("client.go: Invalid json file: %s\n", jsonErr)
		os.Exit(1)
	}

	//1. communicate to DS to get the list of tor nodes
	tnMap := contactDsSerer(clientConfig.DSIp, clientConfig.MaxNumNodes, clientConfig.DSPublicKey)

	if uint16(len(tnMap)) < clientConfig.MaxNumNodes {
		panic("DS did not send enough TN nodes")
	}

	nodeOrder := DetermineTnOrder(tnMap)
	nodeOrder = append(nodeOrder, clientConfig.DSIp)


	//2. create and send onion
	onionMessage, symmKeys := CreateOnionMessage(nodeOrder, tnMap, os.Args[2])

	res := sendOnionMessage(nodeOrder[0], onionMessage, symmKeys)

	fmt.Println("we have received this value from the server: ", res)

}