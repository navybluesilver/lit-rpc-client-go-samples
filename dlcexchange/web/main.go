package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	config "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/config"
	trader "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/trader"

)

const (
	coinType    uint32 = 1
	mHost       string = "localhost"
	mPort       int32  = 8001
	mListenPort uint32 = 2448
	tHost       string = "localhost"
	tPort       int32  = 8002
	tListenPort uint32 = 2449
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
var m *trader.Trader
var t *trader.Trader

type OrderbookPage struct {
	Instrument string
	Underlying string
	MarketMaker string
	MarketMakerLNAddress string
	Oracle string
	Rpoint string
	SPOT int
	SettlementDate string
	Bids interface{}
	Asks interface{}
}

func main() {
	m, _ = trader.NewTrader("Market Maker", mHost, mPort, nil)
	t, _ = trader.NewTrader("Current User", tHost, tPort, m)
	connect(t, m)

	//orderbook
	http.HandleFunc("/", orderbookHandler)
	http.HandleFunc("/buy", buyHandler)
	http.HandleFunc("/sell", sellHandler)

	//files
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))

	//listen
	// redirect every http request to https
	go http.ListenAndServe(port, http.HandlerFunc(redirect))
	log.Fatal(http.ListenAndServeTLS(port, certFile, keyFile, nil))
}

// Connect to the market maker
func connect(t *trader.Trader, m *trader.Trader) {
	mLNAddr, err := m.Lit.GetLNAddress()
	handleError(err)
	fmt.Printf("[%s]- Connecting %s to %s [%s@%s:%d]\n", time.Now().Format("20060102150405"), t.Name, m.Name, mLNAddr, mHost, mListenPort)
	err = t.Lit.Connect(mLNAddr, mHost, mListenPort)
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
	o.Instrument = m.GetInstrument()
	o.Underlying = m.GetUnderlying()
	o.MarketMaker = m.Name
	o.Oracle = m.GetOraclePubKey()
	o.Rpoint = trader.GetR(m.GetSettlementTime())
  o.MarketMakerLNAddress, _ = m.Lit.GetLNAddress()
	o.Bids = m.GetBids()
	o.Asks = m.GetAsks()
	o.SPOT = m.GetCurrentSpot()
	o.SettlementDate = fmt.Sprintf("%v",time.Unix(int64(m.GetSettlementTime()), 0))

	t := template.Must(template.New("orderbook.html").Funcs(fmap).ParseFiles("template/orderbook.html"))
	err := t.ExecuteTemplate(w, "orderbook.html", o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sellHandler(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println("SELL")
 		fmt.Printf("%+v\n", r.Form)
		orderbookHandler(w, r)
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println("BUY")
 		fmt.Printf("%+v\n", r.Form)
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
