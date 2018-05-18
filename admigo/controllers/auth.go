package controllers

import (
	"admigo/common"
	"admigo/model/users"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

const (
	UUID_COOKI string = "uuidcookie"
)

// POST /signup
// Create an user account
func SignupAccount(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	res := users.UserCreate(request)
	output, _ := json.MarshalIndent(res, "", "\t")
	if res.Errors != nil {
		w.WriteHeader(500)
	}
	w.Write(output)
}

// GET /confirm/*filepath
// Confirm registration to user
func ConfirmUser(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	vals := request.URL.Query()
	key := vals.Get("key")
	ckey := common.Encrypt(vals.Get("email"))
	data := map[string]interface{}{"logged": nil, "menuitem": "none"}
	if key != ckey {
		data["msgs"] = []string{
			"Sorry, but You cannot register with this link",
			"Try register again",
		}
		common.GenerateHTML(w, data, "layout", "sidebar", "nav", "error")
		return
	}
	uuid, err := users.UserForceLogin(vals.Get("email"))
	if err != nil {
		data["msgs"] = []string{err.Error()}
		common.GenerateHTML(w, data, "layout", "sidebar", "nav", "error")
		return
	}
	setUuidCookie(w, uuid)
	http.Redirect(w, request, "/", 302)
}

// POST
// Authenticate the user given the email and password
func Login(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	res, uuid := users.UserLogin(request)
	output, _ := json.MarshalIndent(res, "", "\t")
	if res.Errors != nil {
		w.WriteHeader(500)
	}
	setUuidCookie(w, uuid)
	w.Write(output)
}

// User Logout
func Logout(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	user := LoggedUser(request)
	if user == nil {
		return
	}
	err := user.DeleteSessions()
	if err != nil {
		WriteError("all", w, err)
		return
	}
	removeUuidCookie(w)
}

func LoggedUser(r *http.Request) (user *users.UserModel) {
	if cook, err := r.Cookie(UUID_COOKI); err == nil {
		user = users.SessionUser(cook.Value)
	}
	return
}

func setUuidCookie(w http.ResponseWriter, uuid string) {
	cookie := http.Cookie{
		Name:     UUID_COOKI,
		Path:     "/",
		Value:    uuid,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}

func removeUuidCookie(w http.ResponseWriter) {
	rc := http.Cookie{
		Name:    UUID_COOKI,
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, &rc)
}
