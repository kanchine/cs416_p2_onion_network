package dirserver

import (
	"../utils"
	"../keyLibrary"
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	mathrand "math/rand"
	"time"
)

var (
	epochNonce uint64 = 12345
	chCapacity uint8 = 5
	lostMsgThresh uint8 = 5

	Trace = log.New(os.Stdout, "[TRACE] ", 0)
	//Trace = log.New(ioutil.Discard, "[TRACE] ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[ERROR] ", 0)
	//Error = log.New(ioutil.Discard, "[ERROR] ", 0)
)

type DirServer struct {
	Ip     		string
	PortForTN	string
	PortForTC	string
	PortForHB	string
	PriKey 		*rsa.PrivateKey
	TNs    		map[string]rsa.PublicKey
	Fd			utils.FD
	NotifyCh	<-chan utils.FailureDetected
	Mu     		*sync.Mutex
}

func StartDS(Ip, PortForTN, PortForTC, PortForHB string) {

	fmt.Println("==========================================================")
	fmt.Println("Initializing DS...")

	ds := NewDirServer(Ip, PortForTN, PortForTC, PortForHB)
	fmt.Println("DS setup is complete")

	ds.initFD()
	ds.startService()
	ds.startMonitoring()
}

func NewDirServer(Ip, PortForTN, PortForTC, PortForHB string) *DirServer {

	ds := new(DirServer)
	ds.loadPrivateKey()
	ds.TNs = make(map[string]rsa.PublicKey)
	ds.Ip = Ip
	ds.PortForTN = PortForTN
	ds.PortForTC = PortForTC
	ds.PortForHB = PortForHB

	return ds
}

func (ds *DirServer) loadPrivateKey() {

	key, err := keyLibrary.LoadPrivateKey("../dirserver/private.pem")
	checkError(err)

	ds.PriKey = key
}

func (ds *DirServer) initFD() {
	fd, notifyCh, err := utils.Initialize(epochNonce, chCapacity)
	checkError(err)

	ds.Fd = fd
	ds.NotifyCh = notifyCh
}

func (ds *DirServer) startService() {

	go ds.listenAndServeTN()
	go ds.listenAndServeTC()
}

func (ds *DirServer) listenAndServeTN() {

	localTcpAddr, err := net.ResolveTCPAddr("tcp", ds.Ip + ":" + ds.PortForTN)
	checkError(err)

	listener, err := net.ListenTCP("tcp", localTcpAddr)
	checkError(err)

	fmt.Println("Listening on", listener.Addr().String(), "for incoming TNs...")

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			printError("Failed to accept a TN connection request:", err)
			continue
		}
		fmt.Println("================================================================")
		fmt.Println("Here comes a new TN: ", conn.RemoteAddr().String())

		go ds.handleTN(conn)
	}
}

func (ds *DirServer) listenAndServeTC() {

	localTcpAddr, err := net.ResolveTCPAddr("tcp", ds.Ip + ":" + ds.PortForTC)
	checkError(err)

	listener, err := net.ListenTCP("tcp", localTcpAddr)
	checkError(err)

	fmt.Println("Listening on", listener.Addr().String(), "for incoming TCs...")

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			printError("Failed to accept a TC connection request:", err)
			continue
		}
		fmt.Println("================================================================")
		fmt.Println("Here comes a new TC: ", conn.RemoteAddr().String())

		go ds.handleTC(conn)
	}
}

func (ds *DirServer) handleTN(conn *net.TCPConn) {

	defer func() {
		err := conn.Close()
		if err != nil {
			printError("handleTN: failed to close tcp connection.", err)
		}
	}()

	data, err := readFromConnection(conn)
	if err != nil {
		printError("handleTN: reading data from connection failed", err)
		return
	}

	var req utils.NetworkJoinRequest
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		printError("handleTN: request unmarshal failed", err)
		return
	}

	ds.Mu.Lock()
	ds.TNs[req.TorIp] = req.PubKey
	ds.Mu.Unlock()

	var resp utils.NetworkJoinResponse
	resp.Status = true

	err = ds.Fd.AddMonitor(ds.Ip + ":" + ds.PortForHB, req.FdlibIp, lostMsgThresh)
	if err != nil {
		printError("handleTN: AddMonitor failed", err)
		resp.Status = false
	}

	respData, err := json.Marshal(&resp)
	if err != nil {
		printError("handleTN: response marshaling failed", err)
		return
	}

	_, err = writeToConnection(conn, string(respData))
	if err != nil {
		printError("handleTN: response write failed", err)
		return
	}

	if resp.Status {
		Trace.Println("TN: " + req.TorIp + " has joined the Tor network")
		Trace.Println("Start monitoring TN: ", req.FdlibIp)
	}
}

func (ds *DirServer) handleTC(conn *net.TCPConn) {

	defer func() {
		err := conn.Close()
		if err != nil {
			printError("handleTC: failed to close tcp connection.", err)
		}
	}()

	data, err := readFromConnection(conn)
	if err != nil {
		printError("handleTC: reading data from connection failed", err)
		return
	}

	decryptedData, err := keyLibrary.PrivKeyDecrypt(ds.PriKey, []byte(data))
	if err != nil {
		printError("handleTC: request decryption failed", err)
		return
	}

	var req utils.DsRequest
	err = json.Unmarshal(decryptedData, &req)
	if err != nil {
		printError("handleTC: request unmarshal failed", err)
		return
	}

	// Select a specified number of TNs at random. If not enough TNs, return all of them
	circuit := ds.setupCircuit(req.NumNodes)
	var resp utils.DsResponse
	resp.DnMap = circuit

	// Marshall and encrypt the circuit
	respData, err := json.Marshal(&resp)
	if err != nil {
		printError("handleTC: response marshaling failed", err)
		return
	}
	encryptedData, err := keyLibrary.SymmKeyEncrypt(respData, req.SymmKey)
	_, err = writeToConnection(conn, string(encryptedData))
	if err != nil {
		printError("handleTC: response write failed", err)
		return
	}

	Trace.Println("A circuit of ", len(circuit), " TNs has been setup for TC: ", conn.RemoteAddr())
}

func (ds *DirServer) startMonitoring() {

	for {
		select {
		case notify := <-ds.NotifyCh:
			Trace.Println("Detected a failure of", notify)
			ds.removeTN(notify.UDPIpPort)
		case <-time.After(time.Duration(int(lostMsgThresh)*3) * time.Second):
		}
	}
}

func (ds *DirServer) removeTN(TNAddr string) {

	ipToRemove, _, err := net.SplitHostPort(TNAddr)
	if err != nil {
		printError("Failed to get ip of the TN to remove: " + TNAddr, err)
	}

	for addr := range ds.TNs {
		ip, _, err := net.SplitHostPort(addr)
		if err != nil {
			printError("Failed to get ip of the TN to remove: " + TNAddr, err)
			continue
		}

		if ip == ipToRemove {
			ds.Mu.Lock()
			delete(ds.TNs, addr)
			ds.Mu.Unlock()
			Trace.Println("TN: " + addr + " has been removed from Tor network")
		}
	}
}

func (ds *DirServer) setupCircuit(numTNs uint16) map[string]rsa.PublicKey {

	ds.Mu.Lock()
	defer ds.Mu.Unlock()

	if len(ds.TNs) <= int(numTNs) {
		return ds.TNs
	}

	keys := getKeysFromMap(ds.TNs)
	circuit := make(map[string]rsa.PublicKey)

	mathrand.Seed(time.Now().Unix())
	for numTNs > 0 {
		i := mathrand.Intn(len(keys)-1)
		circuit[keys[i]] = ds.TNs[keys[i]]
		keys[i] = keys[len(keys)-1]
		keys = keys[:len(keys)-1]
		numTNs--
	}

	return circuit
}

/**
 *	If there is not a pair of keys available yet, we call this func to generate a pair for use.
 */
func SaveKeysOnDisk() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)

	publicKey := privateKey.PublicKey

	keyLibrary.SavePrivateKeyOnDisk("../dirserver/private.pem", privateKey)
	keyLibrary.SavePublicKeyOnDisk("../dirserver/public.pem", &publicKey)
}

func readFromConnection(conn *net.TCPConn) (string, error) {
	data, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		return "", err
	}

	return data, nil
}

func writeToConnection (conn *net.TCPConn, json string) (int, error) {
	n, werr := fmt.Fprintf(conn, json+"\n")
	return n, werr
}

func getKeysFromMap(m map[string]rsa.PublicKey) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func printError(msg string, err error) {

	Error.Println("****************************************************************")
	Error.Println(msg)
	Error.Println(err)
	Error.Println("****************************************************************")
}

func checkError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}