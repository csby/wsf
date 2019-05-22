package router

import (
	"github.com/csby/wsf/types"
	"net/http"
)

var _ http.Handler = New()

type Router struct {
	trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound func(http.ResponseWriter, *http.Request, types.Assistant)

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})

	Doc types.Doc
}

func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
}

func (s *Router) Document() types.Doc {
	return s.Doc
}

func (s *Router) GET(httpPath types.HttpPath, preHandle types.RouterPreHandle, routerHandle types.RouterHandle, docHandle types.DocHandle) {
	s.Handle("GET", httpPath, preHandle, routerHandle, docHandle)
}

func (s *Router) POST(httpPath types.HttpPath, preHandle types.RouterPreHandle, routerHandle types.RouterHandle, docHandle types.DocHandle) {
	s.Handle("POST", httpPath, preHandle, routerHandle, docHandle)
}

func (s *Router) ServeFiles(httpPath types.HttpPath, preHandle types.RouterPreHandle, root http.FileSystem, docHandle types.DocHandle) {
	if httpPath == nil {
		panic("http path is nil")
	}

	path := httpPath.Path()
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	s.GET(httpPath,
		preHandle,
		func(w http.ResponseWriter, req *http.Request, ps types.Params, _ types.Assistant) {
			req.URL.Path = ps.ByName("filepath")
			fileServer.ServeHTTP(w, req)
		},
		docHandle)
}

func (s *Router) Handle(method string, httpPath types.HttpPath, preHandle types.RouterPreHandle, routerHandle types.RouterHandle, docHandle types.DocHandle) {
	if httpPath == nil {
		panic("http path is nil")
	}

	path := httpPath.Path()
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if s.trees == nil {
		s.trees = make(map[string]*node)
	}

	root := s.trees[method]
	if root == nil {
		root = new(node)
		s.trees[method] = root
	}

	// http
	root.addRoute(path, routerHandle, preHandle)

	// document
	if docHandle != nil {
		if s.Doc != nil {
			if s.Doc.Enable() {
				docHandle(s.Doc, method, httpPath)
			}
		}
	}
}

func (s *Router) Serve(w http.ResponseWriter, req *http.Request, assistant types.Assistant) {
	if s.PanicHandler != nil {
		defer s.recv(w, req)
	}

	path := req.URL.Path

	if root := s.trees[req.Method]; root != nil {
		if handle, preHandle, ps, tsr := root.getValue(path); handle != nil {
			if preHandle != nil {
				if preHandle(w, req, ps, assistant) {
					return
				}
			}
			handle(w, req, ps, assistant)
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && s.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if s.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					s.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	if req.Method == "OPTIONS" && s.HandleOPTIONS {
		// Handle OPTIONS requests
		if allow := s.allowed(path, req.Method); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			return
		}
	} else {
		// Handle 405
		if s.HandleMethodNotAllowed {
			if allow := s.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set("Allow", allow)
				if s.MethodNotAllowed != nil {
					s.MethodNotAllowed.ServeHTTP(w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if s.NotFound != nil {
		s.NotFound(w, req, assistant)
	} else {
		http.NotFound(w, req)
	}
}

func (s *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Serve(w, req, nil)
}

func (s *Router) Lookup(method, path string) (types.RouterHandle, types.RouterPreHandle, Params, bool) {
	if root := s.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, nil, false
}

func (s *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		s.PanicHandler(w, req, rcv)
	}
}

func (s *Router) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range s.trees {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range s.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _, _, _ := s.trees[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}
