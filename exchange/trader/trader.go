package trader

import (
	"fmt"
	"time"
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

type Order struct {
	PeerIdx uint32
	ContractIdx uint64
	AskBidInd string
	Price int
	Quantity int
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

// Buy future for the given price and quantity
// Calculate the funding and division based of the configured margin
func (t *Trader) Buy(price, quantity int) error {
	fmt.Printf("%s offers to buy %d items at %d xBT\n", t.Name, quantity, price)
	ourFunding := int64(price * quantity)
	theirFunding := int64((price * quantity) * margin)
	valueFullyOurs := int64(0)
	valueFullyTheirs := int64(price + (price * margin))
	t.sendContract(ourFunding, theirFunding, valueFullyOurs, valueFullyTheirs)
	return nil
}

// Sell future for the given price and quantity
// Calculate the funding and division based of the configured margin
func (t *Trader) Sell(price, quantity int) error {
	fmt.Printf("%s offers to sell %d items at %d xBT\n", t.Name, quantity, price)
	ourFunding := int64((price * quantity) * margin)
	theirFunding := int64(price * quantity)
	valueFullyOurs := int64(price + (price * margin))
	valueFullyTheirs := int64(0)
	t.sendContract(ourFunding, theirFunding, valueFullyOurs, valueFullyTheirs)
	return nil
}

// Get the balance for a coinType and an address to allow deposits
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

// Accept profitable trades
func (m *Trader) MakeMarket(port uint32) error {
	err := m.Lit.Listen(fmt.Sprintf(":%d", port))
	handleError(err)
	fmt.Printf("Running %s ...\n", m.Name)
	for wait := true; wait; wait = (1 == 1)  {
		m.GetOrderBook()
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}

// Check if the oracle exists
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

// Return the oracle index
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

// Return the market maker peer index
func (t *Trader) getMarketMakerIdx() (uint32, error) {
	return 1, nil
}

// Create and offer the contract to the market maker
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
	fmt.Printf("%s offers contract: ourFunding [%d] | theirFunding [%d] | valueFullyOurs [%d] | valueFullyTheirs [%d]\n", t.Name, ourFunding, theirFunding, valueFullyOurs, valueFullyTheirs)

	return nil
}

func isContractValidOrder(contractIdx int) (b bool, err error) {
	//OurPayoutBase or TheirPayoutBase should be 0, but not both
	//TheirFunding AND OurFunding should not be 0

	//TheirPayoutBase 0 => Sell
	//OurPayoutBase 0 => Buy

	//if a sell
	//OurPayoutBase should be greater than TheirPayoutBase
	//price = OurPayoutBase / (1 + margin)
	//quantity = TheirFunding / price

	//if a Buy
	//TheirPayoutBase should be greater than OurPayoutBase
	//price = TheirPayoutBase / ( 1 + margin)
	//quantity = OurFunding / price
	return true, nil
}

func (t *Trader) convertContractToOrder(contractIdx uint64) (o Order, err error) {
	//TODO: If the contract is not a valid order, decline it

	c, err := t.Lit.GetContract(contractIdx)
	handleError(err)

	o.PeerIdx = c.PeerIdx
	o.ContractIdx = c.Idx


	if c.OurFundingAmount == 0  &&  c.TheirFundingAmount == 0 {
		return o, fmt.Errorf("OurFundingAmount and TheirFundingAmount cannot both be 0")
	}


	// identify if it is a bid or an ask based on the funding
	if c.OurFundingAmount <  c.TheirFundingAmount {
		o.AskBidInd = "BID" // asking a certain price for the instrument, usually higher than the market price, usually triggered by a sell order
	} else {
		o.AskBidInd = "ASK" // bidding a certain price for the instrument, usually lower than the market price, usually triggered by a buy order
	}

	// identify valueFullyOurs and valueFullyTheirs
	var valueFullyOurs int64
	var valueFullyTheirs int64

	if o.AskBidInd == "BID" {
		// if it is a bid, valueFullyOurs is 0
		// valueFullyTheirs is the minimum oracle value that gives us 0
		valueFullyOurs = 0
		for _, d := range c.Division {
			if d.ValueOurs == 0 {
				if valueFullyTheirs == 0 {
					valueFullyTheirs = d.OracleValue
				}
				if d.OracleValue <= valueFullyTheirs  {
	        valueFullyTheirs = d.OracleValue
				}
			}
		}
		o.Price = int(valueFullyTheirs) / (1 + margin)
		o.Quantity = int(c.OurFundingAmount) / o.Price
	} else {
		// if its is an ask, valueFullyTheirs is 0
		// valueFullyOurs is the minimum oracle value that gives us OurFundingAmount + TheirFundingAmount
		valueFullyTheirs = 0
		for _, d := range c.Division {
			if d.ValueOurs == c.OurFundingAmount + c.TheirFundingAmount {
				if valueFullyOurs == 0 {
					valueFullyOurs = d.OracleValue
				}
				if d.OracleValue <= valueFullyOurs  {
					valueFullyOurs = d.OracleValue
				}
			}
		}
		o.Price = int(valueFullyOurs) / (1 + margin)
		o.Quantity = int(c.TheirFundingAmount) / o.Price
	}

	return o, nil
}

//Return all buy and sell offers
func (t *Trader) GetOrderBook() error {
	var orders []Order
	//Get all Contracts
	allContracts, err := t.Lit.ListContracts()
	handleError(err)

	//Loop all contracts offered to t
	//TODO: should not have to loop through all contracts one by oracle
	//TODO: should have the ability to completly remove contracts
	for _, c := range allContracts {
		if c.Status == 2 {
			o, err := t.convertContractToOrder(c.Idx)
			handleError(err)
			orders = append(orders, o)
		}
	}

	fmt.Printf("[%s]: %v\n", time.Now().Format("20060102150405"), orders)

	return nil
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
