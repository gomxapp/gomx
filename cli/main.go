package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
)

var sameDir *bool

func main() {
	sameDir = flag.Bool("same-dir", false, "create project in current directory")
	flag.Parse()

	if flag.Arg(0) == "new" {
		if len(flag.Args()) < 2 {
			fmt.Println("ERROR: Please provide a name")
			fmt.Println("Example: gomx new <app_name>")
			return
		}
		err := newApp(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done!")
	}
}

const templateFileServer = "http://localhost:8081"

func getTemplateFile(filename string) (string, error) {
	templateFileUrl, err := url.JoinPath(templateFileServer, "template", filename)
	if err != nil {
		return "", err
	}
	res, err := http.Get(templateFileUrl)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("request was unsuccessful. Status code: %d", res.StatusCode))
	}
	v := res.Header.Get("Gomx-Version")
	if v == "" {
		return "", errors.New("missing Gomx-Version header. Request was probably unsuccessful")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func createFile(name string, appName string) error {
	templateBody := fmt.Sprintf(`{{define "appName"}}%s{{end}}`, appName)
	fileBody, err := getTemplateFile(name)
	if err != nil {
		return err
	}
	templateBody += fileBody
	println(templateBody)
	bodyT, err := template.New("file").Parse(templateBody)
	if err != nil {
		return err
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	err = bodyT.Execute(file, nil)
	if err != nil {
		return err
	}
	return nil
}

func newApp(appName string) error {
	var rollback = true
	defer func() {
		if rollback {
			fmt.Println("There was an error. Rolling back changes.")
		}
	}()

	var err error
	fmt.Println("Creating new GOMX app with name: " + appName)

	if !(*sameDir) {
		err = os.Mkdir(appName, 0775)
		if err != nil {
			return err
		}
		err = os.Chdir(appName)
		if err != nil {
			return err
		}
		defer func() {
			if rollback {
				err = os.Chdir("..")
				if err != nil {
					log.Fatalln(err)
				}
				err = os.Remove(appName)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}()
	}

	err = createFile("main.go", appName)
	if err != nil {
		return err
	}
	err = createFile("go.mod", appName)
	if err != nil {
		return err
	}
	err = createFile("gomx.config.json", appName)
	if err != nil {
		return err
	}
	err = createFile("tailwind.config.js", appName)
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
	err = createFile("index.gohtml", appName)
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
	err = createFile("input.css", appName)
	if err != nil {
		return err
	}
	rollback = false
	return nil
}
