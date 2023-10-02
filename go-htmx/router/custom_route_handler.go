package router

import (
	"net/http"
	"strings"
)

type route struct {
	path    string
	method  string
	handler http.Handler
}

type CustomRouteHandler struct{ routes []route }

func (h *CustomRouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Static files
	fs := http.FileServer(http.Dir("./client"))
	if strings.HasPrefix(r.URL.Path, "/static/") {
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
		return
	}
	// Pages
	for _, route := range h.routes {
		if route.path == r.URL.Path && route.method == r.Method {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// 404
	handleNotFound(w)
}

func (h *CustomRouteHandler) GET(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  http.MethodGet,
		handler: handler,
	})
}

func (h *CustomRouteHandler) POST(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  http.MethodPost,
		handler: handler,
	})
}

func (h *CustomRouteHandler) PUT(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  http.MethodPut,
		handler: handler,
	})
}

func (h *CustomRouteHandler) DELETE(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  http.MethodDelete,
		handler: handler,
	})
}
