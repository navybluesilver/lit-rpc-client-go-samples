package main

import (
	trader "github.com/mit-dci/lit-rpc-client-go-samples/exchange/trader"
)

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	alice, err := trader.NewTrader("Alice", "127.0.0.1", 8001)
	handleError(err)

	bob, err := trader.NewTrader("Bob", "127.0.0.1", 8002)
	handleError(err)

	//Alice offers to buy
	err = alice.Buy(15000, 100)
	handleError(err)
	//Bob offers to sell
 	err = bob.Sell(15100,100)
	handleError(err)
}
