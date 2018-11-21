package utils

import "crypto/rsa"

type Request struct {
	Key 	string
	SymmKey []byte
}

type Response struct {
	Value 	string
}

type Onion struct {
	NextIp	 string
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

type ClientConfig struct {
	//TODO-ming figure out what else the client needs
	DSPublicKey rsa.PublicKey
	MaxNumNodes uint16
	DSIp		string
	ServerIp	string
}


