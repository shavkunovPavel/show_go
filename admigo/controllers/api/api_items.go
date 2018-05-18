package api

import (
	c "admigo/controllers"
	"admigo/model"
	"admigo/model/items"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func ItemsList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var output []byte
	resp, err := items.List(r)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	output, _ = json.MarshalIndent(resp, "", "\t")
	w.Write(output)
}

func Item(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(ps.ByName("id"))
	item, err := items.ItemById(id)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	output, _ := json.MarshalIndent(item, "", "\t")
	w.Write(output)
}

func EditItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	image, err := c.SaveImage(r, "file", "items")
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	res := items.ItemUpdate(r.FormValue("body"), image)
	output, _ := json.MarshalIndent(res, "", "\t")
	if res.Errors != nil {
		w.WriteHeader(500)
	}
	w.Write(output)
}

func DeleteItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var output []byte
	id, _ := strconv.Atoi(ps.ByName("id"))
	item, err := items.ItemById(id)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	err = item.Delete()
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	ok := model.GetOk("Item was deleted")
	output, _ = json.MarshalIndent(ok, "", "\t")
	w.Write(output)
}
