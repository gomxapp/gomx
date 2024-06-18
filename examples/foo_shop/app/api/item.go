package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/winstonco/gomx/router"
	"gomx.examples.hello_world/data"
)

func init() {
	router.RegisterOnPath("/item/{id}", router.GET,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			arg := r.PathValue("id")
			id, err := strconv.Atoi(arg)
			if err != nil {
				router.ReturnBadRequestSimple(w, err)
				return
			}
			log.Println(id)
			item, err := data.GetItem(id)
			log.Println(item)
			err = router.ReturnGoHTMLFromFiles(w,
				[]string{"../index.gohtml", "item.gohtml"}, "", item)
			if err != nil {
				router.ReturnBadRequestSimple(w, err)
				return
			}
		}),
	)
}
