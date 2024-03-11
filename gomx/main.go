package main

import (
	"fmt"
	// "io"
	// "io/fs"
	"log"
	"os"
	// "path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Commands = []*cli.Command{
		{
			Name:  "greet",
			Usage: "fight the loneliness!",
			Action: func(*cli.Context) error {
				fmt.Println("Hello friend!")
				return nil
			},
		},
		{
			Name:  "new",
			Usage: "create a new gomx app",
			Action: func(ctx *cli.Context) error {
				if ctx.Args().Len() == 0 {
					log.Fatal("Please provide a name")
				}
				err := initGomxApp(ctx.Args().First())
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initGomxApp(target string) error {
	err := os.Mkdir(target, 0775)
	if err != nil {
		return err
	}
	err = os.Chdir(target)
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
`module new_gomx_app

go 1.22

require github.com/winstonco/gomx`)
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
    <title>{{block "title" .}}Default Title{{end}}</title>
    <link rel="stylesheet" href="/static/output.css" />
    <script src="https://unpkg.com/htmx.org@1.9.6"
        integrity="sha384-FhXw7b6AlE/jyjlZH5iHa/tTe9EpJ1Y55RjcgPbjeWMskSxZt1v9qkxLJWNJaGni"
        crossorigin="anonymous"></script>
</head>

<body>
	{{block "body" .}}<h1 class="mb-4 text-4xl font-bold leading-none tracking-tight text-gray-900 md:text-5xl lg:text-6xl dark:text-white">Hello World</h1>{{end}}
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
	// return copyAppHelper(filepath.Join("../examples/", dir), target)
}

// func copyAppHelper(dir string, root string) error {
// 	// log.Println(dir)
// 	// files, err := os.ReadDir(dir)
// 	// log.Println(files)
// 	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
// 		// log.Println(path)
// 		log.Println(path, d.Name())
// 		if err != nil {
// 			return err
// 		}
// 		if d.IsDir() {
// 			log.Println("main.go:78 " + d.Name())
// 			err := os.Mkdir(filepath.Join(root, d.Name()), 0775)
// 			if err != nil {
// 				return err
// 			}
// 			// return copyAppHelper(
// 			// 	filepath.Join(dir, d.Name()),
// 			// 	filepath.Join(root, d.Name()),
// 			// )
// 		}
// 		// fIn, err := os.Open(path)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		// fOut, err := os.Create(root + d.Name())
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		// _, err = io.Copy(fOut, fIn)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		return nil
// 	})
// 	return err
// }
