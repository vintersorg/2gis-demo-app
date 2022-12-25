package app

import (
	"fmt"
	"net/http"

	"github.com/2gis-demo-app/api"
	"github.com/2gis-demo-app/cfg"
	"github.com/2gis-demo-app/log"
	"github.com/2gis-demo-app/oms"
)

type Application struct {
	config cfg.Config
	Logger log.Logger
	api    *api.OrderApi
}

func NewApp() *Application {
	config := cfg.NewConfig()
	logger := log.NewLogger(config)
	provider := oms.NewOMS(logger)
	return &Application{
		config: config,
		Logger: logger,
		api:    api.NewApi(logger, provider),
	}
}

func (a *Application) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", a.config.ListenPort), a.api.Mux)
	if err != nil {
		return err
	}
	return nil
}
