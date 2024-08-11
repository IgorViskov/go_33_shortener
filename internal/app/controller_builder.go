package app

import (
	"net/http"
)

type ControllerBuilder struct {
	controllers []Controller
}

func (cb *ControllerBuilder) AddController(с Controller) {
	cb.controllers = append(cb.controllers, с)
}

func (cb *ControllerBuilder) Build(mux *http.ServeMux, middlewares ...Middleware) {
	for _, controller := range cb.controllers {
		mux.Handle(controller.GetPath(), conveyor(methodHandler(controller), middlewares))
	}
}

func conveyor(h http.Handler, middlewares []Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
