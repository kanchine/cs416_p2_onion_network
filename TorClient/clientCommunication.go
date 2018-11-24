package TorClient

import (
	"../keyLibrary"
	"../utils"
	"bufio"
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

	fmt.Fprint(conn, string(onion) + "\n")

	return readResponse(conn, symmKeys)
}

func readResponse(conn net.Conn, symmKeys [][]byte) string {

	json, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		panic("can not read response from connection")
	}

	return DecryptServerResponse([]byte(json), symmKeys)
}

func sendReqToDs(numNodes uint16, dsPublicKey rsa.PublicKey, conn net.Conn) []byte {
	symmKey := keyLibrary.GenerateSymmKey()

	request := utils.DsRequest{numNodes, symmKey}
	reqBytes, err := utils.Marshall(request)

	if err != nil {
		fmt.Println("Bad marshalling")
	}

	encryptedRequest := EncryptPayload(reqBytes, dsPublicKey)
	marshalledRequest, _ := utils.Marshall(encryptedRequest)

	fmt.Fprintf(conn, string(marshalledRequest) + "\n")

	return symmKey
}

func readResFromDs(conn net.Conn, symmKey []byte) map[string]rsa.PublicKey {
	json, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		panic("can not read response from connection")
	}

	decryptedBytes, err := keyLibrary.SymmKeyDecrypt([]byte(json), symmKey)

	if err != nil {
		panic("can not decrypt response from DS")
	}

	var dsResponse utils.DsResponse
	utils.UnMarshall(decryptedBytes, dsResponse)

	return dsResponse.DnMap
}

func getTcpConnection(ip string) net.Conn {

	tcpConn, err := net.Dial("tcp", ip)
	if err != nil {
		panic("Could not resolve DS IP")
	}

	return tcpConn

}