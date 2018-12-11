#!/bin/bash

cd ~/go/P2-q4d0b-a9h0b-i5g5-v3d0b/
clear
echo "Clearing outdated logs..."
rm *.txt
echo "Running Client to fetch key: '$2'. Using config file '$1'"
go run client/client.go $1 $2