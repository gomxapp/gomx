package router

import (
	"github.com/winstonco/gomx/config"
	"github.com/winstonco/gomx/util"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func Create() *Router {
	if router == nil {
		router = initRouter()
	}
	return router
}

func initRouter() *Router {
	routes := useFileRoutes("")

	h := &Handler{
		routes: routes,
	}
	r := &Router{
		Handler: h,
	}
	return r
}

func useFileRoutes(root string) []route {
	entries, err := os.ReadDir(config.RoutesDir + root)
	if err != nil {
		log.Fatalln(err)
	}

	routes := make([]route, 0)

	files := make([]string, 0)

	path := root
	for _, entry := range entries {
		// Ignore marked files
		if entry.Name()[0] == '_' {
			continue
		}

		if entry.IsDir() {
			fp := path + "/" + entry.Name()
			routes = append(routes, useFileRoutes(fp)...)
		} else {
			if strings.HasSuffix(entry.Name(), ".html") ||
				strings.HasSuffix(entry.Name(), ".gohtml") {
				files = append(files, entry.Name())
			}
		}
	}

	if len(files) > 0 {
		pathNoSlash := strings.Trim(root, "/")
		rootFileIndex := -1
		for i, file := range files {
			if strings.HasPrefix(file, pathNoSlash) {
				if rootFileIndex != -1 {
					log.Fatalf("Error parsing routes\n\tAmbiguous root file. Multiple files named: %s\n", pathNoSlash)
				}
				rootFileIndex = i
			}
			files[i] = filepath.Join(config.RoutesDir, path, file)
		}
		if rootFileIndex == -1 {
			log.Fatalf("Error parsing routes\n\tNo root file found.\n")
		}
		// move "root file" to beginning of slice
		temp := files[0]
		files[0] = files[rootFileIndex]
		files[rootFileIndex] = temp

		// add base template to beginning of slice
		files = append([]string{config.BaseTemplate}, files...)

		t, err := template.New("base").ParseFiles(files...)
		if err != nil {
			log.Fatalf("Error generating template\n\t%v\n", err)
		}

		r := route{
			path:   path,
			method: http.MethodGet,
			handler: &TemplatePageHandler{
				template: t,
			},
		}
		log.Println(r.path)
		routes = append(routes, r)
		// rts = r with trailing-slash
		rts := r
		rts.path = rts.path + "/"
		log.Println(rts.path)
		routes = append(routes, rts)
	}

	return routes
}
