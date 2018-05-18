package api

import (
	c "admigo/controllers"
	"admigo/model/roles"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RolesList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var output []byte
	resp, err := roles.List(r)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	output, _ = json.MarshalIndent(resp, "", "\t")
	w.Write(output)
}
