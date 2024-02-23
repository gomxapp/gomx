package server

import (
	"fmt"
	"github.com/winstonco/gomx/api"
	"github.com/winstonco/gomx/router"
	"log"
	"net/http"
)

func StartServer() {
	var err error = nil

	fmt.Println("-- Creating route handler")
	r := router.Create()
	handler := r.Handler
	fmt.Println("-- Done")

	api.Attach(handler)

	const port = "8080"
	fmt.Println("Listening at http://localhost:" + port)
	err = http.ListenAndServe("localhost:"+port, handler)
	log.Fatal(err)
}
