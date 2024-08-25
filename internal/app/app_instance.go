package app

import "io"

type AppInstance struct {
	services []io.Closer
}

func NewAppInstance() *AppInstance {
	return &AppInstance{
		services: make([]io.Closer, 0),
	}
}

func (app *AppInstance) Close() error {
	errors := make([]error, len(app.services))
	for _, closable := range app.services {
		err := closable.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}
	return nil
}

func (app *AppInstance) AddClosable(c io.Closer) {
	app.services = append(app.services, c)
}
