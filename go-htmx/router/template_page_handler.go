package router

import (
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
		otherRootDir+"/404.html",
		baseTemplatePath,
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
