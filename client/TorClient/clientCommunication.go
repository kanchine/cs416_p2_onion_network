package TorClient

import (
	"crypto/rsa"
	"fmt"
	"net"

	"../../keyLibrary"
	"../../utils"
)

func ContactDsSerer(DSIp string, numNodes uint16, dsPublicKey rsa.PublicKey) (map[string]rsa.PublicKey, error) {

	conn, connErr := getTCPConnection(DSIp)

	if connErr != nil {
		return nil, connErr
	}

	symmKey := sendReqToDs(numNodes, dsPublicKey, conn)

	return readResFromDs(conn, symmKey), nil

}

func SendOnionMessage(t1 string, onion []byte, symmKeys [][]byte) (string, error) {

	conn, connErr := getTCPConnection(t1)

	if connErr != nil {
		return "", connErr
	}

	utils.TCPWrite(conn, onion)

	return readResponse(conn, symmKeys), nil
}

func readResponse(conn *net.TCPConn, symmKeys [][]byte) string {

	bytesRead, err := utils.TCPRead(conn)

	if err != nil {
		panic("can not read response from connection")
	}

	return DecryptServerResponse(bytesRead, symmKeys)
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

func getTCPConnection(ip string) (*net.TCPConn, error) {
	raddr, _ := net.ResolveTCPAddr("tcp", ip)
	return net.DialTCP("tcp", nil, raddr)
}
