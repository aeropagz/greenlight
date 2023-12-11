package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) healthcheckHandler(c echo.Context) error {
	res := "status: avail\n"
	res +=  fmt.Sprintf("enviroment: %s \n", app.config.env)
	res += fmt.Sprintf("version %s \n", version)
	return c.String(http.StatusOK, res)
}
