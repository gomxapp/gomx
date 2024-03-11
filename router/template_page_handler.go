package router

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/winstonco/gomx/config"
)

type PageData struct {
	Arg any
}

type TemplatePageHandler struct {
	template *template.Template
	data     PageData
}

func (tph *TemplatePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tph.template.Execute(w, tph.data)
	if err != nil {
		ieh, err := internalErrorHandler(err)
		if err != nil {
			log.Println(err)
			return
		}
		ieh(w, r)
	}
}

func internalErrorHandler(err error) (http.HandlerFunc, error) {
	_, err2 := os.Stat(config.ReservedDir + "/500.gohtml")
	if err2 != nil {
		log.Println("500.gohtml file not found")
		return nil, err
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		t, err2 := template.ParseFiles(
			config.BaseTemplate,
			config.ReservedDir+"/500.gohtml",
		)
		if err2 != nil {
			log.Println(err2)
			return
		}
		err2 = t.Execute(w, PageData{
			Arg: err,
		})
		if err2 != nil {
			log.Println(err2)
		}
	}
	return http.HandlerFunc(fn), nil
}

func notFoundHandler() (http.HandlerFunc, error) {
	_, err := os.Stat(config.ReservedDir + "/404.gohtml")
	if err != nil {
		log.Println("404.gohtml file not found")
		return nil, err
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		t, err := template.ParseFiles(
			config.BaseTemplate,
			config.ReservedDir+"/404.gohtml",
		)
		if err != nil {
			log.Println(err)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	}
	return http.HandlerFunc(fn), nil
}
