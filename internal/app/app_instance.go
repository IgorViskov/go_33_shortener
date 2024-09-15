package app

import (
	"github.com/IgorViskov/go_33_shortener/internal/ex"
	"io"
)

type Instance struct {
	services []io.Closer
}

func NewAppInstance() *Instance {
	return &Instance{
		services: make([]io.Closer, 0),
	}
}

func (app *Instance) Close() error {
	errors := make([]error, 0)
	for _, closable := range app.services {
		err := closable.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}
	return ex.AggregateErr(errors)
}

func (app *Instance) AddClosable(c io.Closer) {
	app.services = append(app.services, c)
}
