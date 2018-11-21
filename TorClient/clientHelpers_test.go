package TorClient

import (
	"../keyLibrary"
	"../utils"
	"crypto/rsa"
	"fmt"
	"testing"
)

func TestDetermineTnOrder(t *testing.T) {

	myMap := make(map[string]rsa.PublicKey)
	key, _ := keyLibrary.GeneratePrivPubKey()

	myMap["1"] = key.PublicKey
	myMap["2"] = key.PublicKey
	myMap["3"] = key.PublicKey
	myMap["4"] = key.PublicKey
	myMap["5"] = key.PublicKey
	myMap["6"] = key.PublicKey
	myMap["7"] = key.PublicKey
	myMap["8"] = key.PublicKey
	myMap["9"] = key.PublicKey
	myMap["10"] = key.PublicKey
	myMap["11"] = key.PublicKey

	order := determineTnOrder(myMap)
	fmt.Println(order)
}

func TestCreateOnionMessage(t *testing.T) {

	t1Key, _ := keyLibrary.GeneratePrivPubKey()
	t2Key, _ := keyLibrary.GeneratePrivPubKey()
	t3Key, _ := keyLibrary.GeneratePrivPubKey()
	serverKey, _ := keyLibrary.GeneratePrivPubKey()

	myMap := make(map[string]rsa.PublicKey)

	myMap["1"] = t1Key.PublicKey
	myMap["2"] = t2Key.PublicKey
	myMap["3"] = t3Key.PublicKey
	myMap["server"] = serverKey.PublicKey

	order := []string{"1", "2", "3", "server"}

	onion, symmKeys := createOnionMessage(order, myMap, "Hello World")

	fmt.Println(len(symmKeys))

	if onion.NextIp != "2" {
		fmt.Println("should get 2 but received", onion.NextIp)
		t.Log("FAILED")
	}

	if string(onion.SymmKey) != string(symmKeys[0]) {
		fmt.Println("wrong semm key 1")
		t.Log("FAILED")
	}

	decryptedOnion2bytes, _ := keyLibrary.PrivKeyDecrypt(t2Key, onion.Payload)
	var onion2 utils.Onion
	utils.UnMarshall(decryptedOnion2bytes, len(decryptedOnion2bytes), onion2)

	if onion2.NextIp != "3" {
		fmt.Println("should get 3 but received", onion2.NextIp)
		t.Log("FAILED")
	}

	if string(onion2.SymmKey) != string(symmKeys[1]) {
		fmt.Println("wrong semm key 2")
		t.Log("FAILED")
	}

	decryptedOnion3bytes, _ := keyLibrary.PrivKeyDecrypt(t3Key, onion2.Payload)
	var onion3 utils.Onion
	utils.UnMarshall(decryptedOnion3bytes, len(decryptedOnion3bytes), onion3)

	if onion3.NextIp != "server" {
		fmt.Println("should get server but received", onion3.NextIp)
		t.Log("FAILED")
	}

	if string(onion3.SymmKey) != string(symmKeys[2]) {
		fmt.Println("wrong semm key 3")
		t.Log("FAILED")
	}

}

func TestDecryptOnionRes(t *testing.T) {

}