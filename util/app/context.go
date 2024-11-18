package app

import (
	"go-framework/util/app/render"
	"net/http"
)

func (c *Context) JSON(code int, data any) {
	c.render(code, render.Json{Data: data})
}

func (c *Context) render(code int, json render.Render) {
	err := json.Render(c.ResponseWriter)
	if err != nil {
		c.Error()
	}
}

func (c *Context) Error() {
	c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
}
