package trader

import (
	"fmt"
  "github.com/mit-dci/lit-rpc-client-go"
)

type Trader struct {
	Name string
  Lit *litrpcclient.LitRpcClient
}

//Buy future for the given price and quantity
func (t *Trader) Buy(price int, quantity int) error {
  fmt.Printf("%s buys %d items at %d xBT\n", t.Name, quantity, price)
	return nil
}

//Sell future for the given price and quantity
func (t *Trader) Sell(price int, quantity int) error {
  fmt.Printf("%s sells %d items at %d xBT\n", t.Name, quantity, price)
	return nil
}

func NewTrader(name string, host string, port int32) (*Trader, error) {
	t := new(Trader)
  t.Name = name
  l, err := litrpcclient.NewClient(host, port)
  handleError(err)
  t.Lit = l
	return t, nil
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
