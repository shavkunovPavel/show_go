package main

import (
	"admigo/common"
	"fmt"
	"net/http"
	"strings"
)

func getTarget(r *http.Request) (trg *string) {
	e := common.Env()
	for _, red := range *e.Redirects {
		if strings.HasPrefix(r.URL.Path, "/"+red.Prefix) {
			s := fmt.Sprintf("%s://%s:%d%s", red.Protocol, r.Host, red.Port, r.URL.RequestURI())
			trg = &s
			return
		}
	}
	return
}

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}

	httpTarget := getTarget(r)
	if httpTarget != nil {
		http.Redirect(w, r, *httpTarget, http.StatusFound)
		return
	}

	target := "https://" + r.Host + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusFound)
}
