package router

import (
	"github.com/winstonco/gomx/config"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
}

type TemplatePageHandler struct {
	template *template.Template
}

func (tph *TemplatePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := PageData{}
	err := tph.template.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func handleNotFound(w http.ResponseWriter) *TemplatePageHandler {
	w.WriteHeader(http.StatusNotFound)
	t := template.Must(template.New("base").ParseFiles(
		config.BaseTemplate,
		config.ReservedDir+"/404.gohtml",
	))
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &TemplatePageHandler{
		template: t,
	}
}
