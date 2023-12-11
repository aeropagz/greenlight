package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) healthcheckHandler(c echo.Context) error {
	data := map[string]string{
		"status": "available",
		"enviroment": app.config.env,
		"version": version,
	}
	
	return c.JSON(http.StatusOK, data)
}
