package main

import "github.com/labstack/echo/v4"

func (app *application) routes() *echo.Echo {
	e := echo.New()

	e.GET("/v1/healthcheck", app.healthcheckHandler)
	e.POST("/v1/movies", app.createMovieHandler)
	e.GET("/v1/movies/:id", app.showMovieHandler)

	return e
}
