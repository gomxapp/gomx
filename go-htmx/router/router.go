package router

import (
	"fmt"
	"go-htmx/data"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// const clientRootDir string = "./client"
const pagesRootDir string = "./views"
const otherRootDir string = pagesRootDir + "/_reserved"
const baseTemplatePath string = "./client/index.gohtml"

func MakeHandler(isDev bool) http.Handler {
	h := &CustomRouteHandler{
		routes: useFileRoutes("", isDev),
	}

	data.SeedData()

	h.GET("/items", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items := data.GetItems()
		fmt.Println(items)
		t := template.Must(template.ParseFiles("./views/items/_items.gohtml"))
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, items)
	}))

	h.POST("/items", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		r.ParseForm()
		name := r.FormValue("name")
		price, err := strconv.ParseFloat(r.FormValue("price"), 32)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		newItem := data.Item{
			Name:  name,
			Price: float32(price),
		}

		err = data.AddItem(newItem)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		items := data.GetItems()
		item := items[len(items)-1]
		fmt.Println(item)
		t := template.Must(template.ParseFiles("./views/items/_items.gohtml"))
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "item", item)
	}))

	return h
}

func useFileRoutes(root string, isDev bool) []route {
	entries, err := os.ReadDir(pagesRootDir + root)
	if err != nil {
		log.Fatal(err)
	}

	routes := make([]route, 0)

	for _, entry := range entries {
		path := root
		filepath := path + "/" + entry.Name()

		// Ignore marked files
		if entry.Name()[0] == '_' {
			continue
		}

		if entry.IsDir() {
			routes = append(routes, useFileRoutes(filepath, isDev)...)
		} else if strings.HasSuffix(entry.Name(), ".html") ||
			strings.HasSuffix(entry.Name(), ".gohtml") {
			t := template.Must(template.New("base").ParseFiles(
				pagesRootDir+filepath,
				baseTemplatePath,
			))
			// Check for duplicate routes
			for _, route := range routes {
				if route.path == path {
					log.Fatalf("Error parsing routes\n\tFound two files with same path: %s\n", path)
				}
			}
			r := route{
				path:   path,
				method: http.MethodGet,
				handler: &TemplatePageHandler{
					pageTitle: "Go-HTMX",
					template:  t,
				},
			}
			if isDev {
				fmt.Println(r.path)
			}
			routes = append(routes, r)
			// rts = r-trailing-slash
			rts := r
			rts.path = rts.path + "/"
			if isDev {
				fmt.Println(rts.path)
			}
			routes = append(routes, rts)
		}
	}
	return routes
}
