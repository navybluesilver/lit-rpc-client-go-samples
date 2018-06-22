package trader

import (
	"fmt"
  "github.com/mit-dci/lit-rpc-client-go"
)

type Trader struct {
	Name string
  Lit *litrpcclient.LitRpcClient
}

// Return a new trader and connect it to the market maker if delivered
func NewTrader(name string, host string, port int32, m *Trader) (*Trader, error) {
  // Create new trader
  t := new(Trader)
  t.Name = name
  l, err := litrpcclient.NewClient(host, port)
  handleError(err)
  t.Lit = l
	return t, nil
}

//Buy future for the given price and quantity
func (t *Trader) Buy(price int, quantity int) error {
  fmt.Printf("%s offers to buy %d items at %d xBT\n", t.Name, quantity, price)
	return nil
}

//Sell future for the given price and quantity
func (t *Trader) Sell(price int, quantity int) error {
  fmt.Printf("%s offers to sell %d items at %d xBT\n", t.Name, quantity, price)
	return nil
}


func (t *Trader) GetBalance(coinType uint32) {
  addr, err := t.Lit.GetAddresses(coinType,0,false)
  handleError(err)

  allBal, err := t.Lit.ListBalances()
  handleError(err)

  for _, b := range allBal {    
    if b.CoinType == coinType {
        fmt.Printf("Trader: %s | CoinType: %d | SyncHeight: %d | Utxos: %d | WitConf: %d| Channel: %d | Address: %s\n",
          t.Name, b.CoinType, b.SyncHeight, b.TxoTotal, b.MatureWitty, b.ChanTotal, addr)
    }
  }
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//Accept profitable trades
func (m *Trader) MakeMarket(port uint32) error {
  err := m.Lit.Listen(fmt.Sprintf(":%d",port))
  handleError(err)
  fmt.Printf("Running %s ...\n", m.Name)
	return nil
}
