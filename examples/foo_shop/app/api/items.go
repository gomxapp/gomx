package api

import (
	"log"

	"github.com/winstonco/gomx/router"
	"gomx.examples.hello_world/data"

	"net/http"
	"strconv"
)

func init() {
	router.Register(router.GET, func(w http.ResponseWriter, r *http.Request) {
		items := data.GetItems()
		log.Println(items)
		err := router.ReturnGoHTMLFromFiles(w,
			[]string{"items.gohtml"}, "items", items)
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}
	})

	router.Register(router.POST, func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}
		name := r.FormValue("name")
		price, err := strconv.ParseFloat(r.FormValue("price"), 32)
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}

		err = data.AddItem(name, float32(price))
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}

		items := data.GetItems()
		item := items[len(items)-1]
		log.Println(item)

		err = router.ReturnGoHTMLFromFiles(w,
			[]string{"new-item.gohtml"}, "new-item", item)
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}
	})
}
