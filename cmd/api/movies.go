package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (app *application) createMovieHandler(c echo.Context) error {
	return c.String(http.StatusOK, "created movie")
}

func (app *application) showMovieHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Movie no found")
	}

	return c.String(http.StatusOK, fmt.Sprintf("show movie with id: %d", id))
}