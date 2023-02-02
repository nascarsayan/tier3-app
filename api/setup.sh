#!/usr/bin/env bash
apt update
apt install -y curl
curl -OJL https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
rm go1.19.3.linux-amd64.tar.gz
export PATH="$PATH:/usr/local/go/bin"
go mod download
go build -o fruits
