package main

import (
	"github.com/winstonco/gomx/router"
	"github.com/winstonco/gomx/server"
	_ "gomx.examples.hello_world/app/api"
	"log"
	"net/http"
)

func main() {
	r := router.DefaultRouter()
	r.AddStaticFiles("static")
	s := server.NewServer(&http.Server{
		Addr: "localhost:8080",
	}, r)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
