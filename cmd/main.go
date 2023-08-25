package main

import (
	"Sala-2/internal/interface/rest"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Hello, World!")

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	rest.RegisterRoutes(e)

	if err := e.Start(fmt.Sprintf(":5000")); err != nil {
		panic(err)
	}
}
