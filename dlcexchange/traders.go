package main

import (
	"fmt"
	trader "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/trader"
	"math/rand"
	"time"
)

const (
	coinType    uint32 = 1
	mHost       string = "127.0.0.1"
	mPort       int32  = 8001
	mListenPort uint32 = 2448
)

func main() {
	for wait := true; wait; wait = (1 == 1) {
		tradingSimulation()
		time.Sleep(3000000 * time.Millisecond)
	}
}

func tradingSimulation() {
	m, err := trader.NewTrader("Market Maker", mHost, mPort, nil)
	handleError(err)

	fmt.Printf("[%s]- Starting Alice...\n", time.Now().Format("20060102150405"))
	alice, err := trader.NewTrader("Alice", "127.0.0.1", 8002, m)
	connect(alice, m)

	fmt.Printf("[%s]- Starting Bob...\n", time.Now().Format("20060102150405"))
	bob, err := trader.NewTrader("Bob", "127.0.0.1", 8003, m)
	handleError(err)
	connect(bob, m)

	//Alice offers to buy
	err = alice.Buy(randInt(17100, 18000), 1)
	handleError(err)
	alice.SettleExpired()

	//Bob offers to sell
	err = bob.Sell(randInt(16000, 16900), 1)
	handleError(err)
	bob.SettleExpired()
}

// Connect to the market maker
func connect(t *trader.Trader, m *trader.Trader) {
	mLNAddr, err := m.Lit.GetLNAddress()
	handleError(err)
	fmt.Printf("[%s]- Connecting %s to %s [%s@%s:%d]\n", time.Now().Format("20060102150405"), t.Name, m.Name, mLNAddr, mHost, mListenPort)
	err = t.Lit.Connect(mLNAddr, mHost, mListenPort)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func invalidOrders() {
	//
}
