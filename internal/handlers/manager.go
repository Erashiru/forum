package handlers

import (
	"forum/app"
	"forum/internal/config"
	"forum/internal/service"
)

type handlers struct {
	service service.ServiceI
	app     *app.Application
	conf    *config.Config
}

func New(s service.ServiceI, app *app.Application, conf *config.Config) *handlers {
	return &handlers{
		s,
		app,
		conf,
	}
}
