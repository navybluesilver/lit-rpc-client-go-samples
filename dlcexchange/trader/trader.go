package trader

import (
	"encoding/json"
	"fmt"
	"github.com/mit-dci/lit-rpc-client-go"
	orderbook "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/orderbook"
	"net/http"
	"time"
)

const (
	oracleUrl      string = "https://oracle.gertjaap.org"
	oracleName     string = "SPOT"
	datasourceId   uint64 = 2 // xBT/EUR SPOT
	settlementTime uint64 = 1531785600 // 17500 a76c8b4f6fe5770afffb0ad51adf0702d3b666ceee118eca6fbf27e0da9e6024
	coinType       uint32 = 1
	margin         int    = 2
)

type Trader struct {
	Name string
	Lit  *litrpcclient.LitRpcClient
}

type Rpoint struct {
	R string `json:"R"`
}

type OracleSignature struct {
	Signature string `json:"signature"`
	Value     int    `json:"value"`
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
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
	fmt.Printf("[%s]- %s offers to buy %d items at %d xBT\n", time.Now().Format("20060102150405"), t.Name, quantity, price)
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
	fmt.Printf("[%s]- %s offers to sell %d items at %d xBT\n", time.Now().Format("20060102150405"), t.Name, quantity, price)
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
			fmt.Printf("[%s]- Trader: %s | CoinType: %d | SyncHeight: %d | Utxos: %d | WitConf: %d| Channel: %d | Address: %s\n", time.Now().Format("20060102150405"), t.Name, b.CoinType, b.SyncHeight, b.TxoTotal, b.MatureWitty, b.ChanTotal, addr)
		}
	}
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
// TODO: should not be hardcoded to 1
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
	fmt.Printf("[%s]- %s offers contract: ourFunding [%d] | theirFunding [%d] | valueFullyOurs [%d] | valueFullyTheirs [%d]\n", time.Now().Format("20060102150405"), t.Name, ourFunding, theirFunding, valueFullyOurs, valueFullyTheirs)

	return nil
}

// Converts a lit contract into an ask or bid order
func (t *Trader) convertContractToOrder(contractIdx uint64) (o orderbook.Order, err error) {
	// Get the contract from lit and copy the peerIdx and ContractIdx to the order
	c, err := t.Lit.GetContract(contractIdx)
	handleError(err)
	o.PeerIdx = int(c.PeerIdx)
	o.ContractIdx = int(c.Idx)

	// Make sure that both parties provide funding
	if c.OurFundingAmount == 0 && c.TheirFundingAmount == 0 {
		return o, fmt.Errorf("OurFundingAmount and TheirFundingAmount cannot both be 0")
	}

	// identify if it is a bid or an ask based on the funding
	if c.OurFundingAmount < c.TheirFundingAmount {
		o.AskBidInd = "ASK" // asking a certain price for the instrument, usually higher than the market price, usually triggered by a sell order
	} else {
		o.AskBidInd = "BID" // bidding a certain price for the instrument, usually lower than the market price, usually triggered by a buy order
	}

	// identify valueFullyOurs and valueFullyTheirs
	var valueFullyOurs int64
	var valueFullyTheirs int64

	// valueFullyTheirs is the minimum oracle value that gives us 0
	for _, d := range c.Division {
		if d.ValueOurs == 0 {
			if valueFullyTheirs == 0 {
				valueFullyTheirs = d.OracleValue
			}
			if d.OracleValue <= valueFullyTheirs {
				valueFullyTheirs = d.OracleValue
			}
		}
	}

	// valueFullyOurs is the minimum oracle value that gives us OurFundingAmount + TheirFundingAmount
	for _, d := range c.Division {
		if d.ValueOurs == c.OurFundingAmount+c.TheirFundingAmount {
			if valueFullyOurs == 0 {
				valueFullyOurs = d.OracleValue
			}
			if d.OracleValue <= valueFullyOurs {
				valueFullyOurs = d.OracleValue
			}
		}
	}

	if o.AskBidInd == "ASK" {
		// if it is a ask,
		// valueFullyOurs should be 0
		// valueFullyTheirs should not be 0
		if valueFullyOurs != 0 {
			return o, fmt.Errorf("valueFullyOurs for a ask should be 0")
		}
		if valueFullyTheirs == 0 {
			return o, fmt.Errorf("valueFullyTheirs for a ask should not be 0")
		}
		o.Price = int(valueFullyTheirs) / (1 + margin)
		o.Quantity = int(c.OurFundingAmount) / o.Price
	} else {
		// if its is an bid
		// valueFullyOurs should not be 0
		// valueFullyTheirs should be 0
		if valueFullyOurs == 0 {
			return o, fmt.Errorf("valueFullyOurs for a bid should not be 0")
		}

		if valueFullyTheirs != 0 {
			return o, fmt.Errorf("valueFullyTheirs for an bid should be 0")
		}
		o.Price = int(valueFullyOurs) / (1 + margin)
		o.Quantity = int(c.TheirFundingAmount) / o.Price
	}

	return o, nil
}

//Return all buy and sell offers
func (t *Trader) getAllOrders() (orders []orderbook.Order, err error) {
	//Get all Contracts
	allContracts, err := t.Lit.ListContracts()
	handleError(err)

	//Loop all contracts offered to t
	//TODO: should not have to loop through all contracts one by oracle
	//TODO: should have the ability to completly remove contracts
	for _, c := range allContracts {
		// contracts offered to me
		if c.Status == 2 {
			o, err := t.convertContractToOrder(c.Idx)
			if err != nil {
				//TODO: If the contract is not a valid order, decline it
				t.Lit.DeclineContract(c.Idx)
				fmt.Printf("Declined contract [%v]: %v\n", c.Idx, err)
			} else {
				orders = append(orders, o)
			}
		}
	}
	return orders, nil
}

// Accept profitable trades
func (m *Trader) MakeMarket(port uint32) error {
	fmt.Printf("[%s]- Running: %v\n", time.Now().Format("20060102150405"), m.Name)

	// Listen for new offers
	err := m.Lit.Listen(fmt.Sprintf(":%d", port))
	handleError(err)

	// Endless loop to check if there are any contracts to accept
	for loop := true; loop; loop = (1 == 1) {
		allOrders, err := m.getAllOrders()
		handleError(err)
		c, err := orderbook.GetContractsToAccept(allOrders)
		if err != nil {
			if err.Error() == "Nothing to accept." {
				fmt.Printf("[%s]- %v\n", time.Now().Format("20060102150405"), err.Error())
				time.Sleep(10000 * time.Millisecond)
			} else {
				handleError(err)
			}
		} else {
			// TODO: handle not enough satoshis to accept contract
			for _, i := range c {
				err := m.Lit.AcceptContract(uint64(i))
				handleError(err)
				fmt.Printf("[%s]- Accepted contract [%v]\n", time.Now().Format("20060102150405"), i)
			}
		}

		m.SettleExpired()
	}
	return nil
}

// Return Bids
func (m *Trader) GetBids() (bids []orderbook.Order) {
	allOrders, err := m.getAllOrders()
	handleError(err)
	bids, _, _, _, err = orderbook.GetBidsAsks(allOrders)
	handleError(err)
	return bids
}

// Return Asks
func (m *Trader) GetAsks() (asks []orderbook.Order) {
	allOrders, err := m.getAllOrders()
	handleError(err)
	_, asks, _, _, err = orderbook.GetBidsAsks(allOrders)
	handleError(err)
	return asks
}

func GetR(timestamp int) string {
	var r Rpoint

	req, err := http.NewRequest("GET", fmt.Sprintf("https://oracle.gertjaap.org/api/rpoint/2/%d", timestamp), nil)
	handleError(err)

	client := &http.Client{}

	resp, err := client.Do(req)
	handleError(err)

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	handleError(err)

	return r.R
}

func GetOracleSignature() (oracleValue int64, oracleSignature []byte) {
	var sig OracleSignature

	req, err := http.NewRequest("GET", fmt.Sprintf("https://oracle.gertjaap.org/api/publication/%s", GetR(int(settlementTime))), nil)
	handleError(err)

	client := &http.Client{}

	resp, err := client.Do(req)
	handleError(err)

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&sig)
	handleError(err)

	return int64(sig.Value), []byte(sig.Signature)
}

// Settle all the contracts that are past the settlement date
func (t *Trader) SettleExpired() {
	//Get all Contracts
	allContracts, err := t.Lit.ListContracts()
	handleError(err)

	for _, c := range allContracts {
		// contracts active, past settlement date
		if c.Status == 6 && int(c.OracleTimestamp) < int(time.Now().Unix()) {
			v, s := GetOracleSignature()
			t.Lit.SettleContract(c.Idx, v, s)
			fmt.Printf("Settle contract [%v] at %v satoshis\n", c.Idx, v)
		}
	}
}
