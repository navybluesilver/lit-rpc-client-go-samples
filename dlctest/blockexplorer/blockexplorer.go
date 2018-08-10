package blockexplorer

import (
    "fmt"
    "net/http"
    "encoding/json"
  )

const url = "https://testnet.blockexplorer.com"

type Addr struct {
	AddrStr                 string   `json:"addrStr"`
	BalanceSat              int      `json:"balanceSat"`
	UnconfirmedBalanceSat   int      `json:"unconfirmedBalanceSat"`
}

type BlockHeight struct {
	Info struct {
		Blocks          int     `json:"blocks"`
	} `json:"info"`
}



func GetBalance(address string, confirmed bool) (balance int, err error) {
  var addr Addr

  req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/addr/%s", url, address), nil)
  if err != nil {
    return 0, err
  }

  client := &http.Client{}

  resp, err := client.Do(req)
  if err != nil {
    return 0, err
  }

  defer resp.Body.Close()

  if err := json.NewDecoder(resp.Body).Decode(&addr); err != nil {
    return 0, err
  }

  if confirmed {
      return addr.BalanceSat, nil
  } else {
      return addr.UnconfirmedBalanceSat, nil
  }
}

func HasBalance(address string, confirmed bool) (hasBalance bool) {
  bal, err := GetBalance(address, confirmed)

  if err != nil {
    return false
  }

  if bal > 0 {
    return true
  } else {
    return false
  }
}

func GetBlockHeight() (int) {
  var height BlockHeight

  req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/status?q=getBlockHeight", url), nil)
  if err != nil {
    panic(err)
  }

  client := &http.Client{}

  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }

  defer resp.Body.Close()

  if err := json.NewDecoder(resp.Body).Decode(&height); err != nil {
    panic(err)
  }

  return height.Info.Blocks
}
