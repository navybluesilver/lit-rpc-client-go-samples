package main

import (
	trader "github.com/mit-dci/lit-rpc-client-go-samples/exchange/trader"
)

const (
	coinType uint32 = 257
	mHost string = "127.0.0.1"
  mPort int32 = 8001
	mListenPort uint32 = 2448
)

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	m, err := trader.NewTrader("Market Maker", mHost, mPort, nil)
	handleError(err)
	m.MakeMarket(mListenPort)
}
