package gomx

import (
	"fmt"
	"github.com/gomxapp/gomx/config"
	"github.com/gomxapp/gomx/internal"
	"net/http"
	"path"
)

// Router wraps the http.ServeMux. It matches routes using a RouteTree instance
// which allows for more fine-grained error handling. If you are making your own
// router, make sure to call router.Init() before passing it to the server.
type Router struct {
	Mux *http.ServeMux
	// RouteTree is the root of the router's route tree,
	// which matches incoming requests. Registered APIs
	// are added to the route tree.
	routeTree  *internal.RouteTreeWrapper
	routeMaker internal.RouteMaker

	initialized bool
}

// DefaultRouter initializes and returns a Router with default settings.
func DefaultRouter() *Router {
	r := &Router{
		Mux:        http.NewServeMux(),
		routeMaker: internal.FileBasedRouteMaker(),
		routeTree:  nil,
	}
	r.Init()
	return r
}

// Init is required for all routers.
func (router *Router) Init() {
	fmt.Println("-- Initializing router")
	router.routeTree = &internal.RouteTreeWrapper{
		Tree: router.routeMaker.GetRouteTree(),
	}
	router.initApi()
	router.Mux.HandleFunc("/", router.routeTree.ServeNotFound)
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
	router.Mux.Handle("GET /"+dir+"/", http.StripPrefix("/"+dir+"/", fs))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.routeTree.ContainsExactMatch(r) {
		router.routeTree.ServeClosestMatch(w, r)
		return
	}
	router.Mux.ServeHTTP(w, r)
}
