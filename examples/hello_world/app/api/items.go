package api

import (
	"fmt"
	"github.com/winstonco/gomx/api"
	"github.com/winstonco/gomx/config"
	"github.com/winstonco/gomx/router"
	"gomx.examples.hello_world/data"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

func init() {
	data.SeedData()

	api.Register(router.GET, func(w http.ResponseWriter, r *http.Request) {
		items := data.GetItems()
		fmt.Println(items)
		api.ReturnHTML(w, "items.gohtml", items)
	})

	api.Register(router.POST, func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			api.ReturnBadRequestSimple(w, err)
			return
		}
		name := r.FormValue("name")
		price, err := strconv.ParseFloat(r.FormValue("price"), 32)
		if err != nil {
			api.ReturnBadRequestSimple(w, err)
			return
		}
		newItem := data.Item{
			Name:  name,
			Price: float32(price),
		}

		err = data.AddItem(newItem)
		if err != nil {
			api.ReturnBadRequestSimple(w, err)
			return
		}

		items := data.GetItems()
		item := items[len(items)-1]
		fmt.Println(item)
		t, err := template.ParseFiles(filepath.Join(config.ApiRootDir, "items.gohtml"))
		if err != nil {
			api.ReturnBadRequestSimple(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		err = t.ExecuteTemplate(w, "item", item)
		if err != nil {
			api.ReturnBadRequestSimple(w, err)
		}
	})
}
