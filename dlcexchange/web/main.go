package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"strconv"
	"time"
	config "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/config"
	trader "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/trader"
	counterparty "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/counterparty"

)

const (
	coinType    uint32 = 1
	mName 			string = "navybluesilver.net"
	mURL        string = "navybluesilver.net"
	mLNAddress  string = "ln1cgy9632gya8tqfscs7p9grnq0gvu8rwkvd9v0n"
	mHost       string = "128.199.173.181"
	mPort 			uint32 = 2448
	tName  			string = "Alice"
	tHost       string = "localhost"
	tPort       int32  = 8001
	tListenPort uint32 = 2448
)


var (
	templates        = template.Must(template.ParseFiles("template/orderbook.html"))
	certFile         = config.GetString("web.certFile")
	keyFile          = config.GetString("web.keyFile")
	fmap             = template.FuncMap{
		"formatAsSatoshi": formatAsSatoshi,
	}
	port             = config.GetString("web.port")

)
var c *counterparty.Counterparty
var t *trader.Trader

type OrderbookPage struct {
	Instrument string
	Underlying string
	TraderName string
	CounterpartyName string
	CounterpartyLNAddress string
	Oracle string
	Rpoint string
	SPOT int
	SettlementDate string
	Bids interface{}
	Asks interface{}
}

func main() {
	c = &counterparty.Counterparty{Name: mName, LNAddress: mLNAddress, IP: mHost, Port: mPort, URL: mURL }

	t, _ = trader.NewTrader(tName, tHost, tPort)
	connect(t)

	//orderbook
	http.HandleFunc("/", orderbookHandler)
	http.HandleFunc("/buy", buyHandler)
	http.HandleFunc("/sell", sellHandler)

	//files
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))

	//listen
	// redirect every http request to https
	go http.ListenAndServe(port, http.HandlerFunc(redirect))
	log.Fatal(http.ListenAndServe(port, nil))
}

// Connect to the market maker
func connect(t *trader.Trader) {
	fmt.Printf("[%s]- Connecting %s to %s [%s@%s:%d]\n", time.Now().Format("20060102150405"), t.Name, c.Name, c.LNAddress,  c.IP,  c.Port)
	err := t.Lit.Connect( c.LNAddress,  c.IP,  c.Port)
	handleError(err)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	host := strings.Split(req.Host, ":")[0]
	target := "https://" + host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
	http.StatusTemporaryRedirect)
}

//orderbook
func orderbookHandler(w http.ResponseWriter, r *http.Request) {
	var o OrderbookPage
	o.Instrument = t.GetInstrument()
	o.Underlying = t.GetUnderlying()
	o.TraderName = t.Name
	o.Oracle = t.GetOraclePubKey()
	o.Rpoint = trader.GetR(trader.GetSettlementTime())
  o.CounterpartyLNAddress = c.LNAddress
	o.CounterpartyName = c.Name
	o.Bids = t.GetBids()
	o.Asks = t.GetAsks()
	o.SPOT = t.GetCurrentSpot()
	o.SettlementDate = fmt.Sprintf("%v",time.Unix(int64(trader.GetSettlementTime()), 0))

	t := template.Must(template.New("orderbook.html").Funcs(fmap).ParseFiles("template/orderbook.html"))
	err := t.ExecuteTemplate(w, "orderbook.html", o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sellHandler(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		price, _ := strconv.Atoi(fmt.Sprintf("%s", r.Form["price"][0]))
		quantity, _ := strconv.Atoi(fmt.Sprintf("%s", r.Form["quantity"][0]))
		t.Sell(price, quantity)
		orderbookHandler(w, r)
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		price, _ := strconv.Atoi(fmt.Sprintf("%s", r.Form["price"][0]))
		quantity, _ := strconv.Atoi(fmt.Sprintf("%s", r.Form["quantity"][0]))
		t.Buy(price, quantity)
		orderbookHandler(w, r)
}

//formatting
func formatAsSatoshi(satoshi float64) (string, error) {
	if satoshi == 0 {
		return "", nil
	}
	return fmt.Sprintf("%.0f", satoshi), nil
}

//error handling
func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
