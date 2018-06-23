package trader

import (
	"fmt"
	"github.com/mit-dci/lit-rpc-client-go"
)

const (
	oracleUrl      string = "https://oracle.gertjaap.org"
	oracleName     string = "SPOT"
	datasourceId   uint64 = 2 // xBT/EUR SPOT
	settlementTime uint64 = 1528848000
	coinType       uint32 = 257
	margin         int    = 2
)

type Trader struct {
	Name string
	Lit  *litrpcclient.LitRpcClient
}

// Return a new trader
// Import the oracle if it does not yet exist
func NewTrader(name string, host string, port int32, m *Trader) (*Trader, error) {
	t := new(Trader)
	t.Name = name
	l, err := litrpcclient.NewClient(host, port)
	handleError(err)
	t.Lit = l

	exists, err := t.oracleExists()
	if exists != true {
		t.Lit.ImportOracle(oracleUrl, oracleName)
	}
	return t, nil
}

//Buy future for the given price and quantity
//Calculate the funding and division based of the configured margin
func (t *Trader) Buy(price, quantity int) error {
	fmt.Printf("%s offers to buy %d items at %d xBT\n", t.Name, quantity, price)
	ourFunding := int64(price * quantity)
	theirFunding := int64((price * quantity) * margin)
	valueFullyOurs := int64(0)
	valueFullyTheirs := int64(price + (price * margin))
	t.sendContract(ourFunding, theirFunding, valueFullyOurs, valueFullyTheirs)
	return nil
}

//Sell future for the given price and quantity
//Calculate the funding and division based of the configured margin
func (t *Trader) Sell(price, quantity int) error {
	fmt.Printf("%s offers to sell %d items at %d xBT\n", t.Name, quantity, price)
	ourFunding := int64((price * quantity) * margin)
	theirFunding := int64(price * quantity)
	valueFullyOurs := int64(price + (price * margin))
	valueFullyTheirs := int64(0)
	t.sendContract(ourFunding, theirFunding, valueFullyOurs, valueFullyTheirs)
	return nil
}

//Get the balance for a coinType and an address to allow deposits
func (t *Trader) GetBalance(coinType uint32) {
	addr, err := t.Lit.GetAddresses(coinType, 0, false)
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

//Accept profitable trades
func (m *Trader) MakeMarket(port uint32) error {
	err := m.Lit.Listen(fmt.Sprintf(":%d", port))
	handleError(err)
	fmt.Printf("Running %s ...\n", m.Name)
	return nil
}

//Check if the oracle exists
func (t *Trader) oracleExists() (bool, error) {
	allOracles, err := t.Lit.ListOracles()
	handleError(err)

	for _, o := range allOracles {
		if o.Name == oracleName { //TODO: check based on pubkey instead
			return true, nil
		}
	}
	return false, nil
}

//Return the oracle index
func (t *Trader) getOracleIdx(oracleName string) (uint64, error) {
	allOracles, err := t.Lit.ListOracles()
	handleError(err)

	for _, o := range allOracles {
		if o.Name == oracleName {
			return o.Idx, nil
		}
	}
	return 0, fmt.Errorf("Oracle [%s] not found", oracleName)
}

//Return the market maker peer index
func (t *Trader) getMarketMakerIdx() (uint32, error) {
	return 1, nil
}

//Create and offer the contract to the market maker
func (t *Trader) sendContract(ourFunding, theirFunding, valueFullyOurs, valueFullyTheirs int64) error {
	// Create a new empty draft contract
	contract, err := t.Lit.NewContract()
	handleError(err)

	// Get oracle
	oracleIdx, err := t.getOracleIdx(oracleName)
	handleError(err)

	// Configure the contract to use the oracle we need
	err = t.Lit.SetContractOracle(contract.Idx, oracleIdx)
	handleError(err)

	// Set the settlement time
	err = t.Lit.SetContractSettlementTime(contract.Idx, settlementTime)
	handleError(err)

	// Set the coin type of the contract
	err = t.Lit.SetContractCoinType(contract.Idx, coinType)
	handleError(err)

	// Configure the contract datafeed
	err = t.Lit.SetContractDatafeed(contract.Idx, datasourceId)
	handleError(err)

	// Set the contract funding to 1 BTC each
	err = t.Lit.SetContractFunding(contract.Idx, ourFunding, theirFunding)
	handleError(err)

	// Configure the contract division so that Alice get all the
	// funds when the value is 45000, and Bob gets
	// all the funds when the value is 1
	err = t.Lit.SetContractDivision(contract.Idx, valueFullyOurs, valueFullyTheirs)
	handleError(err)

	// Offer the contract to the market maker
	peerIdx, err := t.getMarketMakerIdx()
	err = t.Lit.OfferContract(contract.Idx, peerIdx)
	handleError(err)

	return nil
}

//Return all buy and sell offers
func (m *Trader) GetOrderBook() error {
	return nil
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
