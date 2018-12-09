#!/bin/bash

cd ~/go/P2-q4d0b-a9h0b-i5g5-v3d0b/
clear
echo "Clearing outdated logs..."
rm *.log
rm *.txt
echo "Running Tor node 2"
go run tn/main.go 10.0.0.12:8001 10.0.0.15:4001 10.0.0.15:4002 100000
