package api

import (
	"github.com/winstonco/gomx/config"
	"github.com/winstonco/gomx/router"
	"html/template"
	"net/http"
	"path/filepath"
)

func init() {
	register(router.GET, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(filepath.Join(config.ApiRootDir, "abc.gohtml"))
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		err = t.Execute(w, nil)
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}
	}))
}
