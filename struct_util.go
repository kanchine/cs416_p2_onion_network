package P2_q4d0b_a9h0b_i5g5_v3d0b

import "crypto/rsa"

type Request struct {
	Key 	rsa.PublicKey
	SymmKey []byte
}

type Response struct {
	Value 	string
}

type Onion struct {
	Ip	 	 string
	SymmKey  []byte
	Payload  []byte
}

type NetworkJoinRequest struct {
	TorIp	string
	FdlibIp	string
	PubKey 	rsa.PublicKey
}

type NetworkJoinResponse struct {
	Status 	bool
}

type DsRequest struct {
	NumNodes  uint16
	SymmKey	  []byte
}

type DsResponse struct {
	DnMap 	map[string]rsa.PublicKey
}


