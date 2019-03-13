package types

import "net/http"

type DocHandle func(doc Doc, method, path string)
type RouterHandle func(http.ResponseWriter, *http.Request, Params, Assistant)

type Router interface {
	GET(path string, routerHandle RouterHandle, docHandle DocHandle)
	POST(path string, routerHandle RouterHandle, docHandle DocHandle)

	// path must end with "/*filepath",
	// example: ServeFiles("/src/*filepath", http.Dir("/var/www"), nil)
	ServeFiles(path string, root http.FileSystem, docHandle DocHandle)

	// api document
	Document() Doc
}
