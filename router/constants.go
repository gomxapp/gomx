package router

import "net/http"

type Method string

const (
	GET     Method = http.MethodGet
	HEAD    Method = http.MethodHead
	POST    Method = http.MethodPost
	PUT     Method = http.MethodPut
	DELETE  Method = http.MethodDelete
	CONNECT Method = http.MethodConnect
	OPTIONS Method = http.MethodOptions
	TRACE   Method = http.MethodTrace
	PATCH   Method = http.MethodPatch
)
