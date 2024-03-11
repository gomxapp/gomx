package router

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/winstonco/gomx/config"
)

type RequestMethod string

const (
	GET     RequestMethod = http.MethodGet
	HEAD    RequestMethod = http.MethodHead
	POST    RequestMethod = http.MethodPost
	PUT     RequestMethod = http.MethodPut
	DELETE  RequestMethod = http.MethodDelete
	CONNECT RequestMethod = http.MethodConnect
	OPTIONS RequestMethod = http.MethodOptions
	TRACE   RequestMethod = http.MethodTrace
	PATCH   RequestMethod = http.MethodPatch
)

type route struct {
	path    string
	method  RequestMethod
	handler http.Handler
}

func Init(mux *http.ServeMux) {
	log.Println("-- Creating route handler")
	routes := initRouter(mux)
	log.Println("-- Done")
	log.Println("-- Attaching API routes")
	initApi(mux)
	log.Println("-- Done")
	initHandler(mux, routes)
}

func initRouter(mux *http.ServeMux) []route {
	routes := useFileRoutes("")
	return routes
}

func useFileRoutes(root string) []route {
	entries, err := os.ReadDir(config.RoutesDir + root)
	if err != nil {
		log.Fatalln(err)
	}

	var routes []route

	var files []string

	for _, entry := range entries {
		// Ignore marked files
		if entry.Name()[0] == '_' {
			continue
		}

		if entry.IsDir() {
			fp := root + "/" + entry.Name()
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
			files[i] = filepath.Join(config.RoutesDir, root, file)
		}
		// if "root file" exists, move it to the beginning of files slice
		if rootFileIndex != -1 {
			temp := files[0]
			files[0] = files[rootFileIndex]
			files[rootFileIndex] = temp
		}

		// add base template to beginning of slice
		files = append([]string{config.BaseTemplate}, files...)

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatalf("Error generating template\n\t%v\n", err)
		}

		r := route{
			path:   root + "/",
			method: http.MethodGet,
			handler: &TemplatePageHandler{
				template: t,
			},
		}
		log.Println(r.path)
		routes = append(routes, r)
	}

	return routes
}

func initHandler(mux *http.ServeMux, routes []route) {
	// Static files
	fs := http.FileServer(http.Dir(config.AppRootDir + "/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Pages
	for _, route := range routes {
		if route.path == "" || route.path == "/" {
			path := fmt.Sprintf("%s %s{$}", string(route.method), route.path)
			mux.Handle(path, route.handler)
			continue
		}
		path := fmt.Sprintf("%s %s", string(route.method), route.path)
		mux.Handle(path, route.handler)
	}

	// 404
	nfh, err := notFoundHandler()
	if err != nil {
		log.Println(err)
		return
	}
	mux.Handle("GET /", nfh)
}
