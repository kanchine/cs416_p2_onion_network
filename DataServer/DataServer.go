package DataServer

import (
	"../keyLibrary"
	"P2-q4d0b-a9h0b-i5g5-v3d0b/utils"
	"bufio"
	"crypto/rand"
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
	size := 1024
	privateKey, err := rsa.GenerateKey(rand.Reader, size)

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

	defer func() {
		err := conn.Close()

		if err != nil {
			fmt.Print("Server handler: failed to close tcp connection.")
		}
	}()


	data, err := readFromConnection(conn)

	if err != nil {
		fmt.Println("Server handler: reading data from connection failed")
		return
	}

	decryptedData, err := keyLibrary.PrivKeyDecrypt(s.Key, []byte(data))

	if err != nil {
		fmt.Println("Server handler: request decryption failed")
		return
	}

	var req utils.Request
	err = json.Unmarshal(decryptedData, &req)

	if err != nil {
		fmt.Println("Server handler: request unmarshal failed")
		return
	}

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

	n, err := writeToConnection(conn, string(encryptedData))

	if err != nil {
		fmt.Println("Server handler: response write failed")
		return
	}

	if n != len(encryptedData) {
		fmt.Println("Server handler: incorrect number of bytes written to the connection")
		return
	}
}

//TODO: this function will be replaced by a utility function latter
func readFromConnection(conn *net.TCPConn) (string, error) {
	data, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		return "", err
	}

	return data, nil
}

//TODO: this function will be replaced by a utility function latter
func writeToConnection (conn *net.TCPConn, json string) (int, error) {
	n, werr := fmt.Fprintf(conn, json+"\n")
	return n, werr
}