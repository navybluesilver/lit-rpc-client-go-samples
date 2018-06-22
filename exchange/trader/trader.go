package trader

import (
	"fmt"
)

type Trader struct {
	Name string
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

func NewTrader(name string) (*Trader, error) {
	client := new(Trader)
  client.Name = name
	return client, nil
}
