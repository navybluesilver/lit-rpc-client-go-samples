#!/bin/bash
systemctl stop lit.bob.service
systemctl disable lit.bob.service
cd $GOPATH/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/bob
rm -rf dlc.db 
rm -rf lit.log 
rm -rf ln.db 
rm -rf privkey.hex 
rm -rf testnet3/
systemctl enable lit.bob.service
echo "./bin/lit --dir=bob -v"
echo "systemctl start lit.bob.service && systemctl status lit.bob.service"
