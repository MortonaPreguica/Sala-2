package rest

import (
	"Sala-2/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/teste", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/upload", func(c echo.Context) error {
		return service.UploadImageHandler(c)
	})
}
