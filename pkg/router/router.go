package router

import (
	"fmt"
	"github.com/gomxapp/gomx/pkg/router/internal"
	"net/http"
	"path"

	"github.com/gomxapp/gomx/internal/config"
)

// Router wraps the http.ServeMux. It matches routes using a RouteTree instance
// which allows for more fine-grained error handling. If you are making your own
// router, make sure to call router.Init() before passing it to the server.
type Router struct {
	mux *http.ServeMux
	// RouteTree is the root of the router's route tree,
	// which matches incoming requests. Registered APIs
	// are added to the route tree.
	routeTree  *internal.RouteTreeWrapper
	RouteMaker internal.RouteMaker

	initialized bool
}

// DefaultRouter initializes and returns a Router with default settings.
func DefaultRouter() *Router {
	r := &Router{
		mux:        http.NewServeMux(),
		RouteMaker: internal.FileBasedRouteMaker(),
		routeTree:  nil,
	}
	r.Init()
	return r
}

// Init is required for all routers.
func (router *Router) Init() {
	fmt.Println("-- Initializing router")
	router.routeTree = &internal.RouteTreeWrapper{
		Tree: router.RouteMaker.GetRouteTree(),
	}
	router.initApi()
	router.mux.HandleFunc("/", router.routeTree.ServeNotFound)
	fmt.Println(router.routeTree)
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
	if router.routeTree.ContainsExactMatch(r) {
		router.routeTree.ServeClosestMatch(w, r)
		return
	}
	router.mux.ServeHTTP(w, r)
}
