# GOMX

GOMX, as you may have guessed, is a portmanteau of Go and HTMX. While apps are intended to be built this way, a GOMX app can do much more.

Take a look at some [examples](https://github.com/winstonco/gomx/tree/main/examples) to see common patterns.

## Requirements

`Go v1.22`

## Getting Started

Install the GOMX CLI
```sh
go install github.com/winstonco/gomx@latest
```
Create a new GOMX app
```sh
gomx new <name_of_gomx_app>
```
Get the package
```sh
cd <name_of_gomx_app>
go get github.com/winstonco/gomx
```

> [!WARNING]
> Make sure your go/bin is on your $PATH

### Tailwind

The base app uses Tailwind for styling, but you aren't required to use it, of course.

```sh
npx tailwindcss -i ./app/static/input.css -o ./app/static/output.css --watch
```

### Air

I recommend the [cosmtrek/air](https://github.com/cosmtrek/air/) package for hot reloading of your app. You can follow their instructions on how to set it up.

## API

In `main.go`, there is a commented import for `.../app/api`. The way APIs are set up in GOMX are in an `api` package in your `app` directory (these locations can be changed in `gomx.config.json`).

Register API routes in the `init` function, then uncomment the import in `main.go` to enable them.

Example from `examples/foo_shop`:
```go
package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/winstonco/gomx/router"
	"gomx.examples.hello_world/data"
)

func init() {
	router.RegisterOnPath("/item/{id}/", router.GET, func(w http.ResponseWriter, r *http.Request) {
		arg := r.PathValue("id")
		id, err := strconv.Atoi(arg)
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}
		log.Println(id)
		item, err := data.GetItem(id)
		log.Println(item)
		err = router.ReturnGoHTMLFromFiles(w,
			[]string{"../index.gohtml", "item.gohtml"}, "", item)
		if err != nil {
			router.ReturnBadRequestSimple(w, err)
			return
		}
	})
}
```

