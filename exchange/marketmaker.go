package main

import (
	marketmaker "github.com/mit-dci/lit-rpc-client-go-samples/exchange/marketmaker"
)

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	m, err := marketmaker.NewMarketMaker("Market Maker")
	handleError(err)
	m.Run()
}
