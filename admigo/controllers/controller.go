package controllers

import (
	"admigo/common"
	"admigo/model"
	"encoding/json"
	"fmt"
	"github.com/olahol/go-imageupload"
	"net/http"
)

func setFrontContent(w http.ResponseWriter, r *http.Request, page string, needLog bool) {
	user := LoggedUser(r)
	data := map[string]interface{}{"logged": user, "menuitem": page}
	if user == nil && needLog {
		common.GenerateHTML(w, data, "layout", "sidebar", "uava", "nav", "notlogged")
		return
	}
	common.GenerateHTML(w, data, "layout", "sidebar", "uava", "nav", page)
}

func WriteError(key string, w http.ResponseWriter, err error) {
	res := model.GetErrorResult(map[string]string{key: err.Error()})
	w.WriteHeader(500)
	output, _ := json.MarshalIndent(res, "", "\t")
	w.Write(output)
}

func SaveImage(r *http.Request, name string, path string) (fileName string, err error) {
	width := 150
	height := 150

	uploaded, info, _ := r.FormFile(name)
	if uploaded == nil {
		return
	}
	contentType := info.Header.Get("Content-Type")

	img, err := imageupload.Process(r, name)
	if err != nil {
		return
	}

	var tmpl string

	switch contentType {
	case "image/png":
		img, err = img.ThumbnailPNG(width, height)
		tmpl = "%s.png"
	case "image/jpeg":
		img, err = img.ThumbnailJPEG(width, height, 100)
		tmpl = "%s.jpg"
	}

	if err != nil {
		return
	}

	nm := fmt.Sprintf(tmpl, common.CreateUUID())
	err = img.Save(fmt.Sprintf("./static/images/%s/%s", path, nm))
	if err != nil {
		return
	}
	fileName = nm
	return
}
