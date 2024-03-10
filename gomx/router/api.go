package router

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
	"github.com/winstonco/gomx/util"
)

var apis []route

func initApi(mux *http.ServeMux) {
	for _, api := range apis {
		log.Println(api.method, api.path)
		path := fmt.Sprintf("%s %s", string(api.method), api.path)
		mux.Handle(path, api.handler)
	}
}

func Register(method RequestMethod, handlerFunc http.HandlerFunc) {
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
		log.Fatalf("Error when creating API for file %s\n", file)
	}
	fileName := filepath.ToSlash(filepath.Base(relPath))
	if !strings.HasSuffix(fileName, ".go") {
		log.Fatalln("Tried to register an API route from outside a Go file")
	}
	path := "/" + strings.TrimSuffix(fileName, ".go")
	converted := strings.ReplaceAll(path, "_", "/")

	apis = append(apis, route{
		path:    converted,
		method:  method,
		handler: handlerFunc,
	})
}

func RegisterOnPath(path string, method RequestMethod, handlerFunc http.HandlerFunc) {
	apis = append(apis, route{
		path:    path,
		method:  method,
		handler: handlerFunc,
	})
}

func ReturnGoHTML(w http.ResponseWriter, htmlString string, data any) error {
	t, err := template.New("tmp").Parse(htmlString)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, nil)
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
		return errors.New("Invalid JSON")
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
		return errors.New("Invalid JSON")
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
