# Bitcoin Core Regtest

Complete Steps 1 and 2 in [lit test setup](https://github.com/mit-dci/lit/blob/master/docs/test-setup.md)

# Lit Alice and Bob
Create working directory:
```
cd $HOME
mkdir dlctutorial
cd dlctutorial
mkdir alice
mkdir bob
```

Download and extract the latest release
Double check if below is the latest one [here](https://github.com/mit-dci/lit/releases)

```
wget https://github.com/mit-dci/lit/releases/download/0.1/lit_v0.1_amd64linux.tar.xz
tar -xvf lit_v0.1_amd64linux.tar.xz
```

Create lit configuration files:
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
