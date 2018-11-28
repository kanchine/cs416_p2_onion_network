package tests

import (
	"../DataServer"
	"../TorClient"
	"../utils"
	"P2-q4d0b-a9h0b-i5g5-v3d0b/keyLibrary"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestServerInit(t *testing.T) {
	server, err := DataServer.Initialize("test.json")

	if err != nil {
		t.Errorf("Server initializetion failed")
	}

	key1 := "a"
	valExpected := "test1"

	if val, ok := server.DataBase[key1]; !ok {
		t.Errorf("Server database key %s not found", key1)
	} else {
		if val != valExpected {
			t.Errorf("Server database value for key %s actual: %s  Expected: %s", key1, val, valExpected)
		}
	}

	key2 := "b"
	valExpected = "test2"

	if val, ok := server.DataBase[key2]; !ok {
		t.Errorf("Server database key %s not found", key2)
	} else {
		if val != valExpected {
			t.Errorf("Server database value for key %s actual: %s  Expected: %s", key1, val, valExpected)
		}
	}

	key3 := "c"
	valExpected = "test3"

	if val, ok := server.DataBase[key3]; !ok {
		t.Errorf("Server database key %s not found", key3)
	} else {
		if val != valExpected {
			t.Errorf("Server database value for key %s actual: %s  Expected: %s", key1, val, valExpected)
		}
	}

	expectedIpPort := "127.0.0.1:8080"

	if server.IpPort != expectedIpPort {
		t.Errorf("Server ip port mismatch.")
	}
}

func TestServerRequest(t *testing.T) {
	server, _ := DataServer.Initialize("test.json")

	go server.StartService()

	time.Sleep(2 * time.Second)

	err := sendAndReceive(server.IpPort, server.Key.PublicKey)

	if err != nil {
		t.Errorf("Unable to receive response from server.")
	}

	fmt.Println(server.IpPort)
}

func sendAndReceive(serverIpPort string, serverPublicKey rsa.PublicKey) error {
	tcpLocalAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8081")

	tcpServerAddr, _ := net.ResolveTCPAddr("tcp", serverIpPort)

	tcpConn, _ := net.DialTCP("tcp", tcpLocalAddr, tcpServerAddr)

	myMap := make(map[string]rsa.PublicKey)

	myMap["server"] = serverPublicKey

	order := []string{"server"}

	onionbytes, symKeys := TorClient.CreateOnionMessage(order, myMap, "a")

	_, _ = utils.WriteToConnection(tcpConn, string(onionbytes))

	str, _ := utils.ReadFromConnection(tcpConn)

	var res []byte
	for _, key := range symKeys {
		res, _ = keyLibrary.SymmKeyDecrypt([]byte(str), key)
	}

	var resp utils.Response

	err := json.Unmarshal(res, &resp)

	return err
}
