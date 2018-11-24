package tests

import (
	"../dirserver"
	"../keyLibrary"
	"log"
	"testing"
)

func TestNewDirServer(t *testing.T) {

	dirserver.SaveKeysOnDisk()

	pubKey, err := keyLibrary.LoadPublicKey("../dirserver/public.pem")
	checkError(err)

	ds := dirserver.NewDirServer("localhost", "8001", "8002", "8003")
	if ds.Ip != "localhost" {
		t.Errorf("Incorrect Ip: " + ds.Ip)
	} else if ds.PortForTN != "8001" {
		t.Errorf("Incorrect PortForTN: " + ds.PortForTN)
	} else if ds.PortForTC != "8002" {
		t.Errorf("Incorrect PortForTC: " + ds.PortForTC)
	} else if ds.PortForHB != "8003" {
		t.Errorf("Incorrect PortForHB: " + ds.PortForHB)
	}

	cipherText, _ := keyLibrary.PubKeyEncrypt(pubKey, []byte("Hello World"))
	decryptedBytes, err := keyLibrary.PrivKeyDecrypt(ds.PriKey, cipherText)
	checkError(err)

	plainText := string(decryptedBytes)
	if plainText != "Hello World" {
		t.Errorf("Unmatched decrypted message: " + plainText)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


