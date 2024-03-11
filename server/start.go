package server

import (
	"log"
	"net/http"

	_ "github.com/winstonco/gomx/config"
	"github.com/winstonco/gomx/router"
)

func Start() {
	mux := http.NewServeMux()

	router.Init(mux)

	const port = ":8080"
	server := &http.Server{
		Addr:    "localhost" + port,
		Handler: mux,
	}
	log.Println("Listening at http://localhost" + port)
	err := server.ListenAndServe()
	log.Fatal(err)
}
