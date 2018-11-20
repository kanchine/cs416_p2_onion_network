package P2_q4d0b_a9h0b_i5g5_v3d0b

import "crypto/rsa"

type request struct {
	key 	rsa.PublicKey
	symmKey []byte
}

type response struct {
	value 	string
}

type onion struct {
	ip	 	 string
	symmKey  []byte
	payload  []byte
}

type networkJoinRequest struct {
	torIp	string
	fdlibIp	string
	pubKey 	rsa.PublicKey
}

type networkJoinResponse struct {
	status 	bool
}

type DsRequest struct {
	numNodes  uint16
	symmKey	  []byte
}

type DsResponse struct {
	tnMap 	map[string]rsa.PublicKey
}


