package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomxapp/gomx/internal/config"
	"github.com/gomxapp/gomx/internal/util"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ApiRegisterFunc = func(tree *Router) (string, Method, http.Handler)

var registerFuncs []ApiRegisterFunc

func (router *Router) initApi() {
	for _, registerFunc := range registerFuncs {
		path, method, handler := registerFunc(router)
		_, err := router.routeTree.Tree.AddRelativeChild(path, method, handler, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// Register accepts an ApiRegisterFunc that returns the new API path, method, and handler.
// The function is called during router Init with that router instance.
func Register(registerFunc ApiRegisterFunc) {
	registerFuncs = append(registerFuncs, registerFunc)
}

// RegisterOnPath calls Register with an ApiRegisterFunc that simply returns
// the given path, method, and handler.
func RegisterOnPath(path string, method Method, handler http.Handler) {
	Register(func(tree *Router) (string, Method, http.Handler) {
		return path, method, handler
	})
}

// RegisterOnFile adds a new API with path = caller filename. Underscore ('_')
// characters are replaced with a slash '/' in the path.
//
// Example paths:
//
// items.go => "/items"
// users_{id}_comments.go => "users/{id}/comments/"
func RegisterOnFile(method Method, handler http.Handler) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatalln("Failed to register an API")
	}
	root, err := filepath.Abs(config.RoutesDir)
	if err != nil {
		log.Fatalln(err)
	}
	relPath, err := filepath.Rel(root, file)
	if err != nil {
		log.Printf("Error when creating API for file %s\n", file)
		log.Fatalln(err)
	}
	fileName := filepath.ToSlash(filepath.Base(relPath))
	if !strings.HasSuffix(fileName, ".go") {
		log.Fatalln("Tried to register an API route from outside a Go file... somehow")
	}
	path := "/" + strings.TrimSuffix(fileName, ".go")
	converted := strings.ReplaceAll(path, "_", "/")

	RegisterOnPath(converted, method, handler)
}

func ReturnGoHTML(w http.ResponseWriter, htmlString string, data any) error {
	t, err := template.New("tmp").Parse(htmlString)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

func ReturnGoHTMLFromFiles(w http.ResponseWriter, files []string, name string, data any) error {
	mappedFiles := util.SliceMap(files, func(file string) string {
		return filepath.Join(config.ApiRootDir, file)
	})
	t, err := template.ParseFiles(mappedFiles...)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html")
	if len(name) > 0 {
		err = t.ExecuteTemplate(w, name, data)
	} else {
		err = t.Execute(w, data)
	}
	if err != nil {
		return err
	}
	return nil
}

func ReturnJSON(w http.ResponseWriter, jsonString string) error {
	if !json.Valid([]byte(jsonString)) {
		return errors.New("invalid JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(jsonString))
	if err != nil {
		return err
	}
	return nil
}

func ReturnJSONFromFile(w http.ResponseWriter, file string) error {
	data, err := os.ReadFile(filepath.Join(config.ApiRootDir, file))
	if err != nil {
		return err
	}
	if !json.Valid([]byte(data)) {
		return errors.New("invalid JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = fmt.Fprint(w, data)
	if err != nil {
		return err
	}
	return nil
}

func ReturnBadRequestSimple(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintf(w, "400 Error: %v", v)
	log.Printf("400 Error: %v\n", v)
}
