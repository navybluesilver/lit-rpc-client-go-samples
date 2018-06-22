# Bitcoin Core Regtest

Complete Steps 1 and 2 in [lit test setup](https://github.com/mit-dci/lit/blob/master/docs/test-setup.md)

# Lit Alice and Bob
Create working directory:
```
cd $HOME
mkdir dlcexchange
cd dlcexchange
mkdir alice
mkdir bob
```

Download and build the latest lit code base
```
git clone https://github.com/mit-dci/lit.git
cd lit && go build
cd cmd/lit-af/ && go build

```

Create lit configuration files:
```
cd $HOME/dlcexchange
```

```
nano alice/lit.conf
```

```
reg=localhost
rpchost=0.0.0.0
rpcport=8001
tracker=http://hubris.media.mit.edu:46580
autoListenPort=:2448
autoReconnect=true
autoReconnectInterval=5
```

```
nano bob/lit.conf
```

```
reg=localhost
rpchost=0.0.0.0
rpcport=8002
tracker=http://hubris.media.mit.edu:46580
autoListenPort=:2449
autoReconnect=true
autoReconnectInterval=5
```

Run Alice
```
./lit/lit -v --dir alice
```

Run Bob
```
./lit/lit -v --dir bob
```
