package TorClient

import (
	"../keyLibrary"
	"../utils"
	"crypto/rsa"
	"fmt"
	"math/rand"
	"time"
)

//returns a list of symmetrical keys from T1 to Tn
//and the onion message
func createOnionMessage(nodeOrder []string, tnMap map[string]rsa.PublicKey, reqKey string) (utils.Onion, [][]byte) {

	var onionMessage utils.Onion
	symKeys := make([][]byte, len(nodeOrder))

	for i := len(nodeOrder) - 2; i > 0; i -- {
		var outerOnionMessage utils.Onion

		symmKey := keyLibrary.GenerateSymmKey()
		outerOnionMessage.SymmKey = symmKey
		symKeys = append([][]byte{symmKey}, symKeys...)

		var marshalledPayload []byte

		if i == len(nodeOrder) - 1 {
			//this is the server
			outerOnionMessage.NextIp = ""

			marshalledPayload, _ = utils.Marshall([]byte(reqKey))

		} else {

			outerOnionMessage.NextIp = nodeOrder[i+1]

			marshalledPayload, _ = utils.Marshall(onionMessage)

		}

		nodePublicKey := tnMap[nodeOrder[i]]
		marshalledPayload, _ = keyLibrary.PubKeyEncrypt(&nodePublicKey, marshalledPayload)
		outerOnionMessage.Payload = marshalledPayload
		onionMessage = outerOnionMessage
	}

	return onionMessage, symKeys
}

func decryptServerResponse(onionBytes []byte, symmKeys [][]byte) string {

	for _, key := range symmKeys {
		decryptedOnionBytes, err := keyLibrary.SymmKeyDecrypt(onionBytes, key)

		if err != nil {
			panic("can not decrypt onion using symmKey")
		}

		var unmarshalledOnion utils.Onion
		err = utils.UnMarshall(decryptedOnionBytes, len(decryptedOnionBytes), unmarshalledOnion)

		if err != nil {
			panic("can not unmarshal onion")
		}

		onionBytes = unmarshalledOnion.Payload
	}

	return string(onionBytes)
}

func determineTnOrder(tnMap map[string]rsa.PublicKey) []string {

	keys := getKeysFromMap(tnMap)

	order := make([]string, len(tnMap))

	rand.Seed(time.Now().Unix())
	for len(keys) > 1 {
		i := rand.Intn(len(keys)-1)
		fmt.Println("i", i)
		fmt.Println("key", keys[i])
		order = append(order, keys[i])

		keys[i] = keys[len(keys)-1]
		keys = keys[:len(keys)-1]
	}

	return append(order, keys[0])
}

func getKeysFromMap(m map[string]rsa.PublicKey) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
