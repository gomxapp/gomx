package internal

import (
	"github.com/gomxapp/gomx/config"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type RouteMaker interface {
	GetRouteTree() *RouteTree
}

type fileBasedRouteMaker struct{}

func FileBasedRouteMaker() RouteMaker {
	return &fileBasedRouteMaker{}
}

func (maker *fileBasedRouteMaker) GetRouteTree() *RouteTree {
	rt := createFileBasedRouteTree()
	return rt
}

func createFileBasedRouteTree() *RouteTree {
	rootNode := createRoot()

	// Walks the directory given by dirPath, creates tree nodes and adds them to the parent
	var helper func(*RouteTree, string)
	helper = func(parent *RouteTree, dirPath string) {
		var pathSegment string
		if dirPath == config.RoutesDir {
			pathSegment = "/"
		} else {
			pathSegment = filepath.Base(dirPath)
		}
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return
		}
		var pageFiles []os.DirEntry
		var subDirs []os.DirEntry
		for _, entry := range entries {
			if entry.Name()[0] == '_' {
				continue
			}
			if entry.IsDir() {
				subDirs = append(subDirs, entry)
			} else {
				if strings.HasSuffix(entry.Name(), ".html") ||
					strings.HasSuffix(entry.Name(), ".gohtml") {
					pageFiles = append(pageFiles, entry)
				}
			}
		}
		currentNode := createNode(pathSegment, http.MethodGet, http.NotFoundHandler(), nil)
		// parse files to make current node
		// if there are files to serve, create a tree node
		rootFileIndex := -1
		var fileFullPaths []string
		var notFoundFilePath string
		i := 0 // using separate index counter since some files are skipped
		for _, file := range pageFiles {
			// notFoundFile found
			if file.Name() == "404.html" || file.Name() == "404.gohtml" {
				notFoundFilePath = filepath.Join(dirPath, file.Name())
				continue
			}
			if file.Name() == pathSegment+".html" || file.Name() == pathSegment+".gohtml" {
				if rootFileIndex != -1 {
					log.Fatalf("Error parsing routes\n\tAmbiguous root file. Multiple files named: %s\n", pathSegment)
				}
				rootFileIndex = i
			}
			fileFullPath := filepath.Join(dirPath, file.Name())
			fileFullPaths = append(fileFullPaths, fileFullPath)
			i++
		}
		// if "root file" exists, move it to the beginning of files slice
		if rootFileIndex != -1 {
			temp := fileFullPaths[0]
			fileFullPaths[0] = fileFullPaths[rootFileIndex]
			fileFullPaths[rootFileIndex] = temp
		}
		// add base template to beginning of slice
		fileFullPaths = append([]string{config.BaseTemplate}, fileFullPaths...)
		// page handler
		templ, err := template.ParseFiles(fileFullPaths...)
		if err != nil {
			log.Fatalf("Error generating page template\n\t%v\n", err)
		}
		currentNode.handler = &TemplateHandler{
			template: templ,
		}
		// not found handler
		if notFoundFilePath != "" {
			notFoundTempl, err := templ.Clone()
			if err != nil {
				log.Fatalf("Error generating error page template\n\t%v\n", err)
			}
			notFoundTempl, err = notFoundTempl.ParseFiles(notFoundFilePath)
			if err != nil {
				log.Fatalf("Error generating error page template\n\t%v\n", err)
			}
			currentNode.notFoundHandler = &TemplateHandler{
				template: notFoundTempl,
			}
		}
		err = parent.AddChild(currentNode)
		if err != nil {
			log.Fatalln(err)
		}

		// parse subdirs
		for _, subDir := range subDirs {
			helper(currentNode, filepath.Join(dirPath, subDir.Name()+"/"))
		}
	}
	helper(rootNode, config.RoutesDir)
	return rootNode
}
