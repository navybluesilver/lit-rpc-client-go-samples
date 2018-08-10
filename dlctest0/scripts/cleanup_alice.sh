#!/bin/bash
systemctl stop lit.alice.service
systemctl disable lit.alice.service
cd $GOPATH/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/alice
rm -rf dlc.db 
rm -rf lit.log 
rm -rf ln.db 
rm -rf privkey.hex 
rm -rf testnet3/
systemctl enable lit.alice.service
echo "./bin/lit --dir=alice -v"
echo "systemctl start lit.alice.service && systemctl status lit.alice.service"
