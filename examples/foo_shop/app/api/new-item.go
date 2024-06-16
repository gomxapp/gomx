package api

import (
	"net/http"

	"github.com/winstonco/gomx/router"
)

func init() {
	router.RegisterOnPath("/new-item", router.DELETE,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := router.ReturnJSON(w, "{}")
			if err != nil {
				router.ReturnBadRequestSimple(w, err)
				return
			}
		}),
	)
}
