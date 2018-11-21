package TorClient

import (
	"../keyLibrary"
	"../utils"
	"crypto/rsa"
	"fmt"
	"net"
)

func contactDsSerer(DSIp string, numNodes uint16, dsPublicKey rsa.PublicKey) map[string]rsa.PublicKey {

	conn := getTcpConnection(DSIp)

	symmKey := sendReqToDs(numNodes, dsPublicKey, conn)

	return readResFromDs(conn, symmKey)

}

func sendOnionMessage(t1 string, onion []byte, symmKeys [][]byte ) string {

	conn := getTcpConnection(t1)

	conn.Write(onion)

	return readResponse(conn, symmKeys)
}

func readResponse(conn net.Conn, symmKeys [][]byte) string {
	var onionBytes []byte

	conn.Read(onionBytes)

	return DecryptServerResponse(onionBytes, symmKeys)
}

func sendReqToDs(numNodes uint16, dsPublicKey rsa.PublicKey, conn net.Conn) []byte {
	symmKey := keyLibrary.GenerateSymmKey()

	request := utils.DsRequest{numNodes, symmKey}
	reqBytes, err := utils.Marshall(request)

	if err != nil {
		fmt.Println("Bad marshalling")
	}

	encryptedRequest, err := keyLibrary.PubKeyEncrypt(&dsPublicKey, reqBytes)

	conn.Write(encryptedRequest)

	return symmKey
}

func readResFromDs(conn net.Conn, symmKey []byte) map[string]rsa.PublicKey {
	var resBytes []byte

	conn.Read(resBytes)

	decryptedBytes, err := keyLibrary.SymmKeyDecrypt(resBytes, symmKey)

	if err != nil {
		panic("can not decrypt response from DS")
	}

	var dsResponse utils.DsResponse
	utils.UnMarshall(decryptedBytes, len(decryptedBytes), dsResponse)

	return dsResponse.DnMap
}

func getTcpConnection(ip string) net.Conn {

	tcpConn, err := net.Dial("tcp", ip)
	if err != nil {
		panic("Could not resolve DS IP")
	}

	return tcpConn

}