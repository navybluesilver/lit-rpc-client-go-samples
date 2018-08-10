#!/bin/bash
cd /etc/systemd/system/
rm -rf lit.bob.service
WorkingDirectory="WorkingDirectory="$GOPATH"/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/"
ExecStart="ExecStart="$GOPATH"/src/github.com/mit-dci/lit-rpc-client-go-samples/dlctest/bin/lit --dir=bob"
echo "[Unit]" >> lit.bob.service
echo "Description=Lit Bob" >> lit.bob.service
echo "" >> lit.bob.service
echo "[Service]" >> lit.bob.service
echo "$WorkingDirectory" >> lit.bob.service
echo "$ExecStart" >> lit.bob.service
echo "" >> lit.bob.service
echo "[Install]" >> lit.bob.service
echo "WantedBy=multi-user.target" >> lit.bob.service
systemctl enable lit.bob.service 
systemctl start lit.bob.service 
systemctl status lit.bob.service 


