#!/bin/bash

cd ~/go/P2-q4d0b-a9h0b-i5g5-v3d0b/
clear
echo "Clearing outdated logs..."
rm *.log
rm *.txt
echo "Running Client to fetch key: '$1'"
go run client/client.go config/client.json $1