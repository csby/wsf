package types

import "net/http"

type DocHandle func(doc Doc, method string, path HttpPath)
type RouterHandle func(http.ResponseWriter, *http.Request, Params, Assistant)
type RouterPreHandle func(http.ResponseWriter, *http.Request, Params, Assistant) bool

type Router interface {
	GET(path HttpPath, preHandle RouterPreHandle, routerHandle RouterHandle, docHandle DocHandle)
	POST(path HttpPath, preHandle RouterPreHandle, routerHandle RouterHandle, docHandle DocHandle)

	// path must end with "/*filepath",
	// example: ServeFiles("/src/*filepath", http.Dir("/var/www"), nil)
	ServeFiles(path HttpPath, preHandle RouterPreHandle, root http.FileSystem, docHandle DocHandle)

	// api document
	Document() Doc
}
