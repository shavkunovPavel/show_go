package api

import (
	c "admigo/controllers"
	"admigo/model/csv"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func EthPrices(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var output []byte
	prices, err := csv.PricesEth("1522540800")
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	output, _ = json.MarshalIndent(prices, "", "\t")
	w.Write(output)
}
