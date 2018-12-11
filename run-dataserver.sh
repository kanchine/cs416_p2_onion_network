#!/bin/bash

cd ~/go/P2-q4d0b-a9h0b-i5g5-v3d0b/
clear
echo "Clearing outdated logs..."
rm *.log
rm *.txt
echo "Running Data Server..."
go run server/server.go config/vm-server.json
