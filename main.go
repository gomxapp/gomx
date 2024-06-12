package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "new" {
		if len(flag.Args()) < 2 {
			fmt.Println("ERROR: Please provide a name")
			fmt.Println("Example: gomx new <app_name>")
			return
		}
		err := initGomxApp(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done!")
	}
}

func initGomxApp(appName string) error {
	fmt.Println("Creating new GOMX app with name: " + appName)

	err := os.Mkdir(appName, 0775)
	if err != nil {
		return err
	}
	err = os.Chdir(appName)
	if err != nil {
		return err
	}
	createFile := func(name, body string) error {
		file, err := os.Create(name)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(body)
		if err != nil {
			return err
		}
		return nil
	}
	err = createFile("main.go",
		`package main

import (
	"github.com/winstonco/gomx/server"
	// _ "new_gomx_app/app/api"
)

func main() {
	server.Start()
}`)
	if err != nil {
		return err
	}
	err = createFile("go.mod",
		`module `+appName+`

go 1.22`)
	if err != nil {
		return err
	}
	err = createFile("gomx.config.json",
		`{
  "appRootDir": "./app",
  "apiRootDir": "./api",
  "routesDir": "/routes",
  "reservedDir": "/_",
  "baseTemplate": "/index.gohtml"
}
`)
	if err != nil {
		return err
	}
	err = createFile("tailwind.config.js",
		`/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{html,gohtml}"],
  theme: {
    extend: {},
  },
  plugins: [],
}
`)
	if err != nil {
		return err
	}
	err = os.Mkdir("app", 0775)
	if err != nil {
		return err
	}
	err = os.Chdir("app")
	if err != nil {
		return err
	}
	err = createFile("index.gohtml",
		`<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{block "title" .}}`+appName+`{{end}}</title>
    <link rel="stylesheet" href="/static/output.css" />
    <script src="https://unpkg.com/htmx.org@1.9.6"
        integrity="sha384-FhXw7b6AlE/jyjlZH5iHa/tTe9EpJ1Y55RjcgPbjeWMskSxZt1v9qkxLJWNJaGni"
        crossorigin="anonymous"></script>
</head>

<body>
	{{block "body" .}}
	<h1 class="mb-4 text-4xl font-bold leading-none tracking-tight text-gray-900 md:text-5xl lg:text-6xl dark:text-white">
		Hello World
	</h1>
	{{end}}
</body>

</html>
`)
	if err != nil {
		return err
	}
	err = os.Mkdir("api", 0775)
	if err != nil {
		return err
	}
	err = os.Mkdir("routes", 0775)
	if err != nil {
		return err
	}
	err = os.Mkdir("static", 0775)
	if err != nil {
		return err
	}
	err = os.Chdir("static")
	if err != nil {
		return err
	}
	err = createFile("input.css",
		`@tailwind base;
@tailwind components;
@tailwind utilities;
`)
	if err != nil {
		return err
	}
	return nil
}
