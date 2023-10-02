package router

import (
	"encoding/json"
	"fmt"
	"go-htmx/data"
	"html/template"
	"log"
	"net/http"
	"os"
)

// const clientRootDir string = "./client"
const pagesRootDir string = "./pages"
const otherRootDir string = pagesRootDir + "/_reserved"
const baseTemplatePath string = "./client/index.html"

func MakeHandler(isDev bool) http.Handler {
	h := &CustomRouteHandler{
		routes: useFileRoutes("", isDev),
	}

	data.SeedData()

	h.GET("/items", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items := data.GetItems()
		fmt.Println(items)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(items)
	}))

	h.POST("/items", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newItem data.Item
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res, err := data.AddItem(newItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
		json.NewEncoder(w).Encode(struct {
			Count int `json:"count"`
		}{
			Count: res,
		})
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
		} else {
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
