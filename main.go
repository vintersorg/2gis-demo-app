package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/2gis-demo-app/app"
)

func main() {
	a := app.NewApp()
	err := a.Run()
	if errors.Is(err, http.ErrServerClosed) {
		a.Logger.LogInfo("api closed")
	} else if err != nil {
		a.Logger.LogErrorf("error listening for api: %s", err)
		os.Exit(1)
	}
}
