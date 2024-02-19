package router

import (
	"net/http"
	"strings"
)

type route struct {
	path    string
	method  RequestMethod
	handler http.Handler
}

type Handler struct{ routes []route }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Static files
	fs := http.FileServer(http.Dir("./app"))
	if strings.HasPrefix(r.URL.Path, "/static/") {
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
		return
	}
	// Pages
	for _, route := range h.routes {
		if route.path == r.URL.Path && string(route.method) == r.Method {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// 404
	handleNotFound(w)
}

type RequestMethod string

const (
	GET     RequestMethod = http.MethodGet
	HEAD    RequestMethod = http.MethodHead
	POST    RequestMethod = http.MethodPost
	PUT     RequestMethod = http.MethodPut
	DELETE  RequestMethod = http.MethodDelete
	CONNECT RequestMethod = http.MethodConnect
	OPTIONS RequestMethod = http.MethodOptions
	TRACE   RequestMethod = http.MethodTrace
	PATCH   RequestMethod = http.MethodPatch
)

func (h *Handler) RegisterHandler(path string, method RequestMethod, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  method,
		handler: handler,
	})
}

func (h *Handler) GET(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  GET,
		handler: handler,
	})
}

func (h *Handler) HEAD(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  HEAD,
		handler: handler,
	})
}

func (h *Handler) POST(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  POST,
		handler: handler,
	})
}

func (h *Handler) PUT(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  PUT,
		handler: handler,
	})
}

func (h *Handler) DELETE(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  DELETE,
		handler: handler,
	})
}

func (h *Handler) CONNECT(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  CONNECT,
		handler: handler,
	})
}

func (h *Handler) OPTIONS(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  OPTIONS,
		handler: handler,
	})
}

func (h *Handler) TRACE(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  TRACE,
		handler: handler,
	})
}

func (h *Handler) PATCH(path string, handler http.Handler) {
	h.routes = append(h.routes, route{
		path:    path,
		method:  PATCH,
		handler: handler,
	})
}
