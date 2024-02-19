package router

import (
	"fmt"
	"go-htmx/config"
	"go-htmx/util"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Router struct {
	_routePaths []string
	Routes      []route
	Handler     *Handler
}

func (r *Router) RoutePaths() []string {
	if r._routePaths != nil {
		return r._routePaths
	}
	mapped := util.SliceMap(r.Routes, func(e route) string {
		return e.path
	})
	r._routePaths = mapped
	return r._routePaths
}

var router *Router

func Create(isDev bool) *Router {
	if router == nil {
		router = initRouter(isDev)
	}
	return router
}

func initRouter(isDev bool) *Router {
	routes := useFileRoutes("", isDev)

	h := &Handler{
		routes: routes,
	}
	r := &Router{
		Handler: h,
	}
	return r
}

func useFileRoutes(root string, isDev bool) []route {
	entries, err := os.ReadDir(config.RoutesDir + root)
	if err != nil {
		log.Fatalln(err)
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
				config.RoutesDir+filepath,
				config.BaseTemplate,
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
			// rts = r with trailing-slash
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
