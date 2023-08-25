package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/teste", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
