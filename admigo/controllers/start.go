package controllers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "home", false)
}

func Users(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "users", true)
}

func Items(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "items", true)
}

func Shopping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "shopping", true)
}

func Ico(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "ico", true)
}

func ChartBTC(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "btc", false)
}

func ChartETH(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "eth", false)
}

func Person(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "lk", true)
}

func App(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	setFrontContent(w, r, "app", false)
}
