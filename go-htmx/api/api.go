package api

import (
	"fmt"
	"go-htmx/config"
	"go-htmx/router"
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
	for _, api := range apis {
		fmt.Println(api.path, api.method)
		h.RegisterHandler(api.path, api.method, api.handler)
	}
}

func register(method router.RequestMethod, handler http.Handler) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatalln("FAIL LMAO")
	}
	path := getApiPath(file)
	fmt.Println(path)

	apis = append(apis, api{
		path:    path,
		method:  method,
		handler: handler,
	})
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

func returnBadRequestSimple(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintf(w, "Error: %v", v)
	fmt.Printf("Error: %v\n", v)
}
