#!/bin/bash
echo "Pulling latest lit from github (github.com/mit-dci/lit)"
cd $GOPATH/src/github.com/mit-dci/lit
git pull
echo "Building binaries"
go build
cd cmd/lit-af
go build
echo "Copy binaries to dlctest (github.com/mit-dci/lit-rpc-client-go-samples/dlctest/bin/)"
cp lit-af $GOPATH/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/bin/lit-af
cd $GOPATH/src/github.com/mit-dci/lit
cp lit $GOPATH/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/bin/lit
cd $GOPATH/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest
