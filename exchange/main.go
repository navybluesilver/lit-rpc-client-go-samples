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
	alice, err := trader.NewTrader("Alice")
	handleError(err)

	bob, err := trader.NewTrader("Bob")
	handleError(err)

	//Alice offers to buy 100 USDU18 at 15000 satoshis
	err = alice.Buy(15000, 100)
	handleError(err)
	//Bob offers to sell 100 USDU18 at 15000 satoshis
 	err = bob.Sell(15000,100)
	handleError(err)
}
