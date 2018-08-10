#!/bin/bash
cd /etc/systemd/system/
rm -rf lit.alice.service
WorkingDirectory="WorkingDirectory="$GOPATH"/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/"
ExecStart="ExecStart="$GOPATH"/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/bin/lit --dir=alice"
echo "[Unit]" >> lit.alice.service
echo "Description=Lit Alice" >> lit.alice.service
echo "" >> lit.alice.service
echo "[Service]" >> lit.alice.service
echo "$WorkingDirectory" >> lit.alice.service
echo "$ExecStart" >> lit.alice.service
echo "" >> lit.alice.service
echo "[Install]" >> lit.alice.service
echo "WantedBy=multi-user.target" >> lit.alice.service
systemctl enable lit.alice.service 
systemctl start lit.alice.service 
systemctl status lit.alice.service 


