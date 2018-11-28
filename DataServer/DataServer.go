package DataServer

import (
	"../keyLibrary"
	"../utils"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
)

const TCP_PROTO = "tcp"

type Server struct {
	Key *rsa.PrivateKey                // Private key of the server, the public key is also in this data structure.
	IpPort string                      // The ip port the server will be listening for connection on.
	DataBase map[string] string        // The key value pair this database stores.
	LockDataBase *sync.Mutex           // Lock to ensure synchronized database access.
}

type Config struct {
	IncomingTcpAddr string             // The ip port the server will be listening for connection on
	DataBase map[string] string        // The key value pair for the data base
}

func Initialize(configFile string) (*Server, error) {

	jsonFile, err := os.Open(configFile)

	defer func() {
		err := jsonFile.Close()

		if err != nil {
			fmt.Print("Server init: json config file closing failed, continue.")
		}
	}()

	if err != nil {
		fmt.Println("Server init: Error opening the configuration file, please try again.")
		return nil, err
	}

	configData, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println("Server init: Error reading the configuration file, please try again")
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configData, &config)

	// TODO: replace this with loading key from the config file
	privateKey, err := keyLibrary.GeneratePrivPubKey()

	return &Server{privateKey, config.IncomingTcpAddr, config.DataBase, &sync.Mutex{}}, err
}

func (s *Server) StartService() {
	localTcpAddr, err := net.ResolveTCPAddr(TCP_PROTO, s.IpPort)

	if err != nil {
		fmt.Println("Listener creation failed, please try again.")
		return
	}

	listener, err := net.ListenTCP(TCP_PROTO, localTcpAddr)


	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("TCP connection failed with client:", tcpConn.RemoteAddr().String())
			continue
		} else {
			fmt.Println("Incoming connection established with client:",tcpConn.RemoteAddr().String())
		}

		go s.connectionHandler(tcpConn)
	}
}

func (s *Server) connectionHandler(conn *net.TCPConn) {
	// Note connection will be closed by the TN.

	reqStr, err := utils.ReadFromConnection(conn)

	if err != nil {
		fmt.Println("Server handler: reading data from connection failed")
		return
	}

	req := unmarshalServerRequest([]byte(reqStr), s.Key)

	var resp utils.Response
	s.LockDataBase.Lock()
	if val, ok := s.DataBase[req.Key]; ok {
		resp.Value = val
	}
	s.LockDataBase.Unlock()

	respData, err := json.Marshal(&resp)
	if err != nil {
		fmt.Println("Server handler: response marshaling failed")
		return
	}

	encryptedData, err := keyLibrary.SymmKeyEncrypt(respData, req.SymmKey)

	n, err := utils.WriteToConnection(conn, string(encryptedData))

	if err != nil {
		fmt.Println("Server handler: response write failed")
		return
	}

	if n != len(encryptedData) + 1 {
		fmt.Println("Server handler: incorrect number of bytes written to the connection")
		return
	}
}

func unmarshalServerRequest(data []byte, serverKey *rsa.PrivateKey) utils.Request{

	var serverBytes [][]byte

	err := utils.UnMarshall(data, &serverBytes)
	if err != nil {
		fmt.Println("Error unmarshal client requests bytes:", err)
	}

	var decryptedServerBytes []byte

	for i := range serverBytes {
		decryptedBytePiece, err := keyLibrary.PrivKeyDecrypt(serverKey, serverBytes[i])
		if err != nil {
			fmt.Println("failed to decrypt client requests:", err)
		}
		decryptedServerBytes = append(decryptedServerBytes, decryptedBytePiece...)
	}

	var serverMessage utils.Request

	err = utils.UnMarshall(decryptedServerBytes, &serverMessage)

	if err != nil {
		fmt.Println("Error unmarshal client requests:", err)
	}

	return serverMessage
}
