package api

import (
	"fmt"
	"go-htmx/config"
	"go-htmx/data"
	"go-htmx/router"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

func init() {
	data.SeedData()

	register(router.GET, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items := data.GetItems()
		fmt.Println(items)
		t, err := template.ParseFiles(filepath.Join(config.ApiRootDir, "items.gohtml"))
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		err = t.Execute(w, items)
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}
	}))

	register(router.POST, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}
		name := r.FormValue("name")
		price, err := strconv.ParseFloat(r.FormValue("price"), 32)
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}
		newItem := data.Item{
			Name:  name,
			Price: float32(price),
		}

		err = data.AddItem(newItem)
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}

		items := data.GetItems()
		item := items[len(items)-1]
		fmt.Println(item)
		t, err := template.ParseFiles(filepath.Join(config.ApiRootDir, "items.gohtml"))
		if err != nil {
			returnBadRequestSimple(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		err = t.ExecuteTemplate(w, "item", item)
		if err != nil {
			returnBadRequestSimple(w, err)
		}
	}))
}
