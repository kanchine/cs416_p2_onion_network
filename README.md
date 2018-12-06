# P2-q4d0b-a9h0b-i5g5-v3d0b
Tor Network - Protect Your Identity!

## How to start Diretory_Server
Run "go run dirserver/dirserver.go [Ip] [PortForTN] [PortForTC]"   

(Default: Ip=localhost, PortForTN=8001, PortForTC=8002)
   
## How to start Data Server
1. Run "go run server/server.go config/server.json"


## How to run tor client
1. go run client/client.go config/client.json keyToFetch

## How to run tor node
1. go run tn/main.go [dsIPPort] [listenIPPort] [fdListenIPPort] [timeOutMillis]

(Default: dsIPPort=127.0.0.1:8001, listenIPPort=127.0.0.1:4001, fdListenIPPort=127.0.0.1:4002, timeOutMillis=1000)
