package router

import (
	"go-htmx/config"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	PageTitle string
}

type TemplatePageHandler struct {
	pageTitle string
	template  *template.Template
}

func (tph *TemplatePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		PageTitle: tph.pageTitle,
	}
	err := tph.template.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func handleNotFound(w http.ResponseWriter) *TemplatePageHandler {
	w.WriteHeader(http.StatusNotFound)
	t := template.Must(template.New("base").ParseFiles(
		config.ReservedDir+"/404.html",
		config.BaseTemplate,
	))
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &TemplatePageHandler{
		pageTitle: "404 Page Not Found",
		template:  t,
	}
}
