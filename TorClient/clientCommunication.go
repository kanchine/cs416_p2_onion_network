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

	encryptedBytes, _ := keyLibrary.PubKeyEncrypt(&dsPublicKey, reqBytes)
	conn.Write(encryptedBytes)

	return symmKey
}

func readResFromDs(conn net.Conn, symmKey []byte) map[string]rsa.PublicKey {
	// TODO: Not sure if it's safe to hardcode the buf size here
	buf := make([]byte, 8192)
	n, err := conn.Read(buf)
	if err != nil {
		panic("can not read response from connection")
	}

	decryptedBytes, err := keyLibrary.SymmKeyDecryptBase64(buf[:n], symmKey)
	if err != nil {
		fmt.Println(err)
		panic("can not decrypt response from DS")
	}

	var dsResponse utils.DsResponse
	err = utils.UnMarshall(decryptedBytes, &dsResponse)
	if err != nil {
		fmt.Println(err)
		panic("readResFromDs: Unmarshalling failed")
	}

	return dsResponse.DnMap
}

func getTcpConnection(ip string) net.Conn {

	tcpConn, err := net.Dial("tcp", ip)
	if err != nil {
		panic("Could not resolve DS IP")
	}

	return tcpConn

}