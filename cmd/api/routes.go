package main

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "greenlight.aeropagz.de/docs/greenlight"
)

func (app *application) routes() *echo.Echo {
	e := echo.New()

	e.GET("/v1/healthcheck", app.healthcheckHandler)
	e.GET("/v1/movies", app.listMovieHandler)
	e.POST("/v1/movies", app.createMovieHandler)
	e.GET("/v1/movies/:id", app.showMovieHandler)
	e.PUT("/v1/movies/:id", app.updateMovieHandler)
	e.DELETE("/v1/movies/:id", app.deleteMovieHandler)

	e.POST("/v1/users", app.createUserHandler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
