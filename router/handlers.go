package router

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Arg any
}

type TemplateHandler struct {
	template *template.Template
	data     PageData
}

func (tph *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tph.template.Execute(w, tph.data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error executing template", 500)
	}
}
