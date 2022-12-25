package app

import (
	"applicationDesignTest/cfg"
	"log"
)

type Application struct {
	Config cfg.Config
	Logger *log.Logger
}

func NewApp() *Application {
	config := cfg.NewConfig()
	return &Application{
		Config: config,
		Logger: log.Default(),
	}
}
