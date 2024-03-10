package api

import (
	"net/http"

	"github.com/winstonco/gomx/router"
)

func init() {
	router.Register(router.DELETE, func(w http.ResponseWriter, r *http.Request) {
		err := router.ReturnJSON(w, "{}")
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}
	})
}
