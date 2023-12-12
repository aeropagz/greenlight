package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"greenlight.aeropagz.de/internal/data"
)

// @Summary     Show server status
// @Tags		Status
// @Accept 		json
// @Produce 	json
// @Success		200 {object} data.Healthcheck
// @Router	/v1/healthcheck [get]
func (app *application) healthcheckHandler(c echo.Context) error {
	data := &data.Healthcheck{
		Status: "up",
		Enviroment: "development",
		Version: version,
	}	
	return c.JSON(http.StatusOK, data)
}
