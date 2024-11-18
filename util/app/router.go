package app

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	router      *httprouter.Router
	middlewares []MiddlewareFunc
}

func NewRouter() *Router {
	return &Router{router: httprouter.New(), middlewares: []MiddlewareFunc{}}
}

func (r *Router) Use(mw MiddlewareFunc) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) addRoute(method string, path string, handle HandlerFunc) {
	finalHandler := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		ctx := &Context{
			ResponseWriter: w,
			Request:        req,
			Params:         ps,
		}
		handle(ctx)
	}
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = adaptMiddleware(r.middlewares[i], finalHandler)
	}
	r.router.Handle(method, path, finalHandler)
}

func adaptMiddleware(mw MiddlewareFunc, handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		newHandle := mw(func(ctx *Context) {
			handle(w, req, ps)
		})
		newHandle(&Context{ResponseWriter: w, Request: req, Params: ps})
	}
}

func (r *Router) GET(path string, handle HandlerFunc) {
	r.addRoute("GET", path, handle)
}

func (r *Router) POST(path string, handle HandlerFunc) {
	r.addRoute("POST", path, handle)
}

func (r *Router) PUT(path string, handle HandlerFunc) {
	r.addRoute("PUT", path, handle)
}

func (r *Router) DELETE(path string, handle HandlerFunc) {
	r.addRoute("DELETE", path, handle)
}

type Group struct {
	prefix      string
	middlewares []MiddlewareFunc
	parent      *Router
}

func (r *Router) Group(prefix string) *Group {
	return &Group{
		prefix: prefix,
		parent: r,
	}
}

func (g *Group) Use(mw MiddlewareFunc) {
	g.middlewares = append(g.middlewares, mw)
}

func (g *Group) addRoute(method, path string, handle HandlerFunc) {
	finalPath := g.prefix + path
	finalHandler := handle

	// Combine group middlewares and router middlewares
	allMiddlewares := append(g.middlewares, g.parent.middlewares...)
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		finalHandler = allMiddlewares[i](finalHandler)
	}

	adaptedHandler := adaptHandler(finalHandler)
	g.parent.router.Handle(method, finalPath, adaptedHandler)
}

// adaptHandler converts a app.HandlerFunc (our custom type) into httprouter.Handle
// which is the type expected by httprouter.
func adaptHandler(h HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		ctx := &Context{
			ResponseWriter: w,
			Request:        req,
			Params:         ps,
		}
		h(ctx)
	}
}

func (g *Group) GET(path string, handle HandlerFunc) {
	g.addRoute("GET", path, handle)
}

func (g *Group) POST(path string, handle HandlerFunc) {
	g.addRoute("POST", path, handle)
}

// Implement ServeHTTP to satisfy http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
