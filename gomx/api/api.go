package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/winstonco/gomx/config"
	"github.com/winstonco/gomx/router"
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
		log.Fatalln("Failed to register an API")
	}
	path := getApiPath(file)

	apis = append(apis, api{
		path:    path,
		method:  method,
		handler: handlerFunc,
	})
}

func ReturnGoHTML(w http.ResponseWriter, htmlString string, data any) error {
	t, err := template.New("temp").Parse(htmlString)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}

func ReturnGoHTMLFromFile(w http.ResponseWriter, file string, data any) error {
	t, err := template.ParseFiles(filepath.Join(config.ApiRootDir, file))
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

func ReturnJSON(w http.ResponseWriter, jsonString string) error {
	if !json.Valid([]byte(jsonString)) {
		return errors.New("Invalid JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
		return errors.New("Invalid JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	fmt.Printf("400 Error: %v\n", v)
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
