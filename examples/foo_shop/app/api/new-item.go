package api

import (
	"github.com/winstonco/gomx/api"
	"github.com/winstonco/gomx/router"
	"html/template"
	"net/http"
)

func init() {
	api.Register(router.DELETE, func(w http.ResponseWriter, r *http.Request) {
		// TODO: ReturnJSON func, ReturnHTMLString func
		t, err := template.New("new-item").Parse("<div id=\"item-added\"></div>")
		if err != nil {
			api.ReturnBadRequestSimple(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		err = t.Execute(w, nil)
		if err != nil {
			api.ReturnBadRequestSimple(w, err)
		}
	})
}
