package api

import (
	"admigo/common"
	c "admigo/controllers"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	api = "api"
)

func UserCanDo(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		logged := c.LoggedUser(r)

		if logged != nil && logged.CanWrite() {
			h(w, r, ps)
			return
		}

		c.WriteError(api, w, errors.New(common.Mess().InsufRights))
	}
}
