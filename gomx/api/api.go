package api

import (
	"fmt"
	"github.com/winstonco/gomx/config"
	"github.com/winstonco/gomx/router"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

type api struct {
	path    string
	method  router.RequestMethod
	handler http.Handler
}

var apis []api

func Attach(h *router.Handler) {
	fmt.Println("-- Attaching API routes")
	for _, api := range apis {
		fmt.Println(api.path, api.method)
		h.RegisterHandler(api.path, api.method, api.handler)
	}
	fmt.Println("-- Done")
}

func Register(method router.RequestMethod, handlerFunc http.HandlerFunc) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatalln("FAIL LMAO")
	}
	path := getApiPath(file)

	apis = append(apis, api{
		path:    path,
		method:  method,
		handler: handlerFunc,
	})
}

func ReturnHTML(w http.ResponseWriter, file string, data any) {
	t, err := template.ParseFiles(filepath.Join(config.ApiRootDir, file))
	if err != nil {
		ReturnBadRequestSimple(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		ReturnBadRequestSimple(w, err)
		return
	}
}

func ReturnBadRequestSimple(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintf(w, "Error: %v", v)
	fmt.Printf("Error: %v\n", v)
}

func getApiPath(fileAbsPath string) string {
	root, err := filepath.Abs(config.RoutesDir)
	if err != nil {
		log.Fatalln(err)
	}
	relPath, err := filepath.Rel(root, fileAbsPath)
	if err != nil {
		log.Fatalf("Error when creating API for file %s\n", fileAbsPath)
	}
	fileName := filepath.ToSlash(filepath.Base(relPath))
	if !strings.HasSuffix(fileName, ".go") {
		log.Fatalln("Tried to register an API route from outside a Go file")
	}
	apiPath := "/" + strings.TrimSuffix(fileName, ".go")
	return apiPath
}
