package common

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

func getTemplates(folder string, filenames []string, fm template.FuncMap) (tmpl *template.Template) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s%s.html", folder, file))
	}
	if fm != nil {
		tmpl = template.Must(template.New("").Funcs(fm).ParseFiles(files...))
		return
	}
	tmpl = template.Must(template.ParseFiles(files...))
	return
}

func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	funcMap := template.FuncMap{"sb": createSidebar}
	getTemplates("", filenames, funcMap).ExecuteTemplate(writer, "layout", data)
}

func GenerateMail(bufer *bytes.Buffer, data interface{}, filenames ...string) {
	getTemplates("mail/", filenames, nil).ExecuteTemplate(bufer, "layout", data)
}

func createSidebar(manutag string) template.HTML {
	return template.HTML(getSidebarHtml(sidebar.Sidebar, manutag))
}
