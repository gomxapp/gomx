package router

import (
	"fmt"
	"net/http"
	"path"

	"github.com/winstonco/gomx/config"
)

// Router wraps the http.ServeMux. It matches routes using a RouteTree instance
// which allows for more fine-grained error handling. If you are making your own
// router, make sure to call router.Init() before passing it to the server.
type Router struct {
	mux *http.ServeMux
	// RouteTree is the root of the router's route tree,
	// which matches incoming requests. Registered APIs
	// are added to the route tree.
	RouteTree  *RouteTreeWrapper
	RouteMaker RouteMaker

	initialized bool
}

// DefaultRouter initializes and returns a Router with default settings.
func DefaultRouter() *Router {
	r := &Router{
		mux:        http.NewServeMux(),
		RouteMaker: FileBasedRouteMaker(),
		RouteTree:  nil,
	}
	r.Init()
	return r
}

// Init is required for all routers.
func (router *Router) Init() {
	fmt.Println("-- Initializing router")
	router.RouteTree = &RouteTreeWrapper{
		Tree: router.RouteMaker.GetRouteTree(),
	}
	router.initApi()
	router.mux.HandleFunc("/", router.RouteTree.ServeNotFound)
	fmt.Println(router.RouteTree)
	fmt.Println("-- Done")
	router.initialized = true
}

func (router *Router) IsInitialized() bool {
	return router.initialized
}

// AddStaticFiles Adds an http.FileServer handler
//
// dir: the directory path following config.AppRootDir
func (router *Router) AddStaticFiles(dir string) {
	// Static files
	fs := http.FileServer(http.Dir(path.Join(config.AppRootDir, dir)))
	router.mux.Handle("GET /"+dir+"/", http.StripPrefix("/"+dir+"/", fs))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtw := router.RouteTree
	p := r.URL.EscapedPath()
	method := Method(r.Method)
	rtw.setClosestMatch(p, method)
	rtw.setPathValues(r)
	if rtw.matchLvl >= wildMatch && rtw.closestNode != nil && rtw.closestNode.handler != nil {
		// if exact match found, just serve with handler
		rtw.closestNode.handler.ServeHTTP(w, r)
		return
	}
	router.mux.ServeHTTP(w, r)
}
