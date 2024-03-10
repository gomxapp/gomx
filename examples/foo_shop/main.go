package main

import (
	"github.com/winstonco/gomx/server"
	_ "gomx.examples.hello_world/app/api"
)

func main() {
	server.Start()
}
