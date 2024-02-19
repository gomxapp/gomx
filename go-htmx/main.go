package main

import (
	"fmt"
	"go-htmx/api"
	"go-htmx/router"
	"log"
	"net/http"
)

func main() {
	startServer(true)
}

func startServer(isDev bool) {
	var err error = nil

	if isDev {
		fmt.Println("Creating route handler")
	}
	r := router.Create(isDev)
	handler := r.Handler
	api.Attach(handler)
	if isDev {
		fmt.Println("Done")
	}

	const port = "8080"
	fmt.Println("Listening at http://localhost:" + port)
	err = http.ListenAndServe("localhost:"+port, handler)
	log.Fatal(err)
}
