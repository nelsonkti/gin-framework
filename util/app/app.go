package app

import "net/http"

type Engine struct {
	*Router
}

func New() *Engine {
	return &Engine{
		NewRouter(),
	}
}

func (a *Engine) Run(addr string) {
	http.ListenAndServe(addr, a)
}
