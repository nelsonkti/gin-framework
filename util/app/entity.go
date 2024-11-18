package app

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
)

type HandlerFunc func(*Context)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	path           string
	pnames         []string
	pvalues        []string
	query          url.Values
	Params         httprouter.Params
}

type Response struct {
	Writer http.ResponseWriter
}

func (c *Context) Path() string {
	return c.path
}

func (c *Context) SetPath(p string) {
	c.path = p
}

func (c *Context) GetRequest() *http.Request {
	return c.Request
}

func (c *Context) Param(name string) string {
	for i, n := range c.pnames {
		if i < len(c.pvalues) {
			if n == name {
				return c.pvalues[i]
			}
		}
	}
	return ""
}

func (c *Context) ParamNames() []string {
	return c.pnames
}
