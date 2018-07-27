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
	port        string = ":80"
	coinType    uint32 = 1
	mHost       string = "127.0.0.1"
	mPort       int32  = 8001
	mListenPort uint32 = 2448
)

var (
	templates        = template.Must(template.ParseFiles("template/orderbook.html"))
	certFile         = config.GetString("web.certFile")
	keyFile          = config.GetString("web.keyFile")
	fmap             = template.FuncMap{
		"formatAsSatoshi": formatAsSatoshi,
	}
)

type OrderbookPage struct {
	Instrument string
	Underlying string
	MarketMakerLNAddress string
	MarketMaker string
	Oracle string
	Rpoint string
	SPOT int
	SettlementDate string
	Bids interface{}
	Asks interface{}
}

func main() {
	//orderbook
	http.HandleFunc("/orderbook", orderbookHandler)

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
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
}

//orderbook
func orderbookHandler(w http.ResponseWriter, r *http.Request) {
	m, err := trader.NewTrader("Market Maker", mHost, mPort, nil)
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
	err = t.ExecuteTemplate(w, "orderbook.html", o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
