package tornode

import (
	"crypto/rsa"
	"net"

	"../../utils"
)

func contactDS(dsIPPort string, TorIPPort string, fdlibIPPort string, pubKey *rsa.PublicKey) (bool, error) {
	var laddr, raddr *net.TCPAddr
	var addrErr error
	laddr, addrErr = net.ResolveTCPAddr("tcp", ":"+getNewUnusedPort())
	raddr, addrErr = net.ResolveTCPAddr("tcp", dsIPPort)
	if addrErr != nil {
		return false, addrErr
	}
	conn, connErr := net.DialTCP("tcp", laddr, raddr)
	if connErr != nil {
		return false, connErr
	}

	request := utils.NetworkJoinRequest{
		TorIpPort:   TorIPPort,
		FdlibIpPort: fdlibIPPort,
		PubKey:      *pubKey,
	}
	payload, merr := utils.Marshall(request)
	if merr != nil {
		return false, merr
	}
	_, werr := utils.WriteToConnection(conn, string(payload))
	if werr != nil {
		return false, werr
	}
	responsePayload, rerr := utils.ReadFromConnection(conn)
	if rerr != nil {
		return false, rerr
	}
	response := &utils.NetworkJoinResponse{}
	umerr := utils.UnMarshall([]byte(responsePayload), response)
	if umerr != nil {
		return false, umerr
	}
	return response.Status, nil
}
