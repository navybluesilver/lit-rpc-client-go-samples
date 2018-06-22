package marketmaker

import (
	"fmt"
  "time"
)

type MarketMaker struct {
	Name string
}

//Accept profitable trades
func (m *MarketMaker) Run() error {
  fmt.Printf("Running %s ...\n", m.Name)
  time.Sleep(5 * time.Second)
	return nil
}


func NewMarketMaker(name string) (*MarketMaker, error) {
	client := new(MarketMaker)
  client.Name = name
	return client, nil
}
