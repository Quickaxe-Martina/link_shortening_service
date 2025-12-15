#!/bin/bash
set -e

export PATH=$PATH:$(go env GOPATH)/bin
goimports -w .
# pkill -9 "shortener"
go build -o ./cmd/loadtest/loadtest ./cmd/loadtest/loadtest.go

chmod +x ./cmd/loadtest/loadtest

./cmd/loadtest/loadtest