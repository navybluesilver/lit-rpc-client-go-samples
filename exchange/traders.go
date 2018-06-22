package main

import (
	"fmt"
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

	fmt.Println("Starting Alice...")
	alice, err := trader.NewTrader("Alice", "127.0.0.1", 8002, m)
	alice.GetBalance(coinType)
	connect(alice, m)
	fmt.Println("Alice Done.")

	fmt.Println("")

	fmt.Println("Starting Bob...")
	bob, err := trader.NewTrader("Bob", "127.0.0.1", 8003, m)
	handleError(err)
	bob.GetBalance(coinType)
	connect(bob, m)
	fmt.Println("Bob Done.")

	fmt.Println("")

	//Alice offers to buy
	err = alice.Buy(15000, 100)
	handleError(err)

	//Bob offers to sell
 	err = bob.Sell(15100,100)
	handleError(err)

	fmt.Println("Done.")
}

// Connect to given market maker if it was delivered
func connect(t *trader.Trader, m *trader.Trader)  {
	mLNAddr, err := m.Lit.GetLNAddress()
	handleError(err)
	fmt.Printf("Connecting %s to %s [%s@%s:%d]\n", t.Name, m.Name, mLNAddr, mHost, mListenPort)
	err = t.Lit.Connect(mLNAddr, mHost, mListenPort)
	handleError(err)
}
