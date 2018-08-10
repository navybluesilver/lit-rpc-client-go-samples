package main

import (
    "fmt"
    "time"
    "os/exec"
    "github.com/mit-dci/lit-rpc-client-go"
    "github.com/mit-dci/lit-rpc-client-go-samples/dlctest/blockexplorer"
  )

type LitConnection struct {
  Name string
  Host string //hostname or ip
  Port int
}

var (
  alice = LitConnection { Name: "Alice", Host: "127.0.0.1", Port: 8001 }
  bob = LitConnection { Name: "Bob", Host: "127.0.0.1", Port: 8002 }
  coinType = 1
)

func main() {

  // make sure we have the latest git
  update_binaries()

  // restart services
  //restartServices() //TODO: format error, do manual for now

  // make sure that the services are running
  ping(alice)
  ping(bob)

  // wait for confirmation
  waitForFunding()

  // create and send contract
  createAndSendContract(alice, bob)

  // accept contract
  acceptContract(bob)

}

// Services
func update_binaries() {
  runScript("./scripts/update_binaries.sh")

}

func restartServices() {
  runScript("./scripts/restart_services.sh")
}

func runScript(script string) {
  cmd := exec.Command(script)
  stdout, err := cmd.Output()

  if err != nil {
      println(err.Error())
      return
  }
  print(string(stdout))
}


// Tasks
func ping(conn LitConnection) {
  lit := getLitClient(conn)
  isListening, err := lit.IsListening()
  handleError(err)


  if isListening {
    log(fmt.Sprintf("[%s] - Listening", conn.Name))
  } else {
    log(fmt.Sprintf("[%s] - Not listening\n", conn.Name))
  }

}

func hasFunding(conn LitConnection) (hasFunding bool) {

  witnessAddress := getWitnessAddress(conn)
  legacyAddress := getLegacyAddress(conn)


  // confirmed witness funding, return true
  if hasBalance(conn, true) {
    return true
  }

  // confirmed legacy funding, sweep Witness Addres, return false and wait
  if blockexplorer.HasBalance(legacyAddress,true) {
        if !checkBlockHeight(conn) {
          return false
        }
        log(fmt.Sprintf("Legacy Address has confirmed funding, now sweeping to the Witness Address: %s", witnessAddress))
        sweepFunds(conn)
        return false
  }

  // unconfirmed legacy funding, return false and wait
  if blockexplorer.HasBalance(legacyAddress,false) {
          log(fmt.Sprintf("Legacy Address has been funded, but still waiting for confirmations on the blockchain: %s", legacyAddress))
          return false
  }


  log(fmt.Sprintf("Legacy Address needs funding: %s", legacyAddress))
  return false
}

func sweepFunds(conn LitConnection) {
  lit := getLitClient(conn)
  lit.Sweep(getWitnessAddress(conn),1)
}

// Requests
func getLNAddress(conn LitConnection) {

}

func getLegacyAddress(conn LitConnection) (pubKey string) {
  lit := getLitClient(conn)
  addr, err := lit.GetAddresses(uint32(coinType), 0, true)
  handleError(err)
  return addr[0]
}

func getWitnessAddress(conn LitConnection) (pubKey string) {
  lit := getLitClient(conn)
  addr, err := lit.GetAddresses(uint32(coinType), 0, false)
  handleError(err)
  return addr[0]
}

func getBalance(conn LitConnection, witness bool) (int) {
  lit := getLitClient(conn)
  bal, err := lit.ListBalances()
  handleError(err)

  utxo := 0
  matureWitty := 0
  for _, b := range bal {
      fmt.Printf("[%s] - Channel: %d | UTXO: %d | Confirmed Witness: %d \n", conn.Name, b.ChanTotal, b.TxoTotal, b.MatureWitty)
      utxo = utxo + int(b.TxoTotal)
      matureWitty = matureWitty + int(b.MatureWitty)
  }

  if witness {
      return matureWitty
  } else {
    return utxo
  }
}

func hasBalance(conn LitConnection, witness bool) (bool) {
  bal := getBalance(conn, witness)
  if bal > 0 {
      return true
  } else {
    return false
  }
}


func checkBlockHeight(conn LitConnection) (ok bool) {
  litHeight := getBlockHeight(conn)
  testnetHeight := blockexplorer.GetBlockHeight()
  delta := testnetHeight - litHeight
  if delta < 0 {
    return true
  }

  if delta > 0 {
    log(fmt.Sprintf("Block height for %s is only [%d], while expecting [%d]", conn.Name, litHeight, testnetHeight))
    return false
  }

  return true
}

func getBlockHeight(conn LitConnection) (int) {
  lit := getLitClient(conn)
  bal, err := lit.ListBalances()
  handleError(err)
  return int(bal[0].SyncHeight)
}

func createAndSendContract(sender LitConnection, receiver LitConnection) {

}

func acceptContract(conn LitConnection) {
  // error if no contract is available to accept
}

func getLitClient(conn LitConnection) (lit *litrpcclient.LitRpcClient) {
  lit, err := litrpcclient.NewClient(conn.Host, int32(conn.Port))
  handleError(err)
  return lit
}

// Waits
func delaySecond(n time.Duration) {
    log("waiting...")
    time.Sleep(n * time.Second)
}

func waitForFunding() {
    for wait := !bothHasFunding(); wait; wait = !bothHasFunding() {
      delaySecond(60)
    }
}

func bothHasFunding() (bool) {
  fundingAlice := hasFunding(alice)
  fundingBob := hasFunding(bob)
  if fundingAlice && fundingBob {
    return true
  }
  return false
}

// Error Handling
func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func log(message string) {
  fmt.Println(message)
}
