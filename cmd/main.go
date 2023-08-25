package main

import (
	"Sala-2/internal/environment"
	"Sala-2/internal/interface/rest"
	"fmt"

	"github.com/bavatech/envloader"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := envloader.Load(&environment.Env); err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	rest.RegisterRoutes(e)

	if err := e.Start(fmt.Sprintf(":%s", environment.Env.ServerPort)); err != nil {
		panic(err)
	}
}
