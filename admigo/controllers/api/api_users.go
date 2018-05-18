package api

import (
	c "admigo/controllers"
	"admigo/model"
	"admigo/model/users"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func UsersList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var output []byte
	resp, err := users.List(r)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	output, _ = json.MarshalIndent(resp, "", "\t")
	w.Write(output)
}

func User(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(ps.ByName("id"))
	user, err := users.UserById(id, true)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	output, _ := json.MarshalIndent(user, "", "\t")
	w.Write(output)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	user, err := users.UserById(id, false)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	logged := c.LoggedUser(r)
	err = user.Delete(logged)
	if err != nil {
		c.WriteError(api, w, err)
		return
	}
	ok := model.GetOk("User was deleted")
	output, _ := json.MarshalIndent(ok, "", "\t")
	w.Write(output)
}

func EditUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	image, err := c.SaveImage(r, "file", "users")
	if err != nil {
		c.WriteError(api, w, err)
		return
	}

	logged := c.LoggedUser(r)
	res := users.UserUpdate(r.FormValue("body"), image, logged)

	output, _ := json.MarshalIndent(res, "", "\t")
	if res.Errors != nil {
		w.WriteHeader(500)
	}

	w.Write(output)
}
