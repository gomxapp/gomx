package main

import (
	"fmt"
	"log"
	"net/http"

	"go-htmx/router"
)

func main() {
	startServer(true)
}

func startServer(isDev bool) {
	var err error = nil

	if isDev {
		fmt.Println("Creating route handler")
	}
	handler := router.MakeHandler(isDev)
	if isDev {
		fmt.Println("Done")
	}

	const port = "8080"
	fmt.Println("Listening at http://localhost:" + port)
	err = http.ListenAndServe("localhost:"+port, handler)
	log.Fatal(err)
}
