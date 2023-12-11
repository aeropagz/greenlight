package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"greenlight.aeropagz.de/internal/data"
)

func (app *application) createMovieHandler(c echo.Context) error {
	return c.String(http.StatusOK, "created movie")
}

// @Summary     Show an a movie
// @Tags		Movie
// @Accept 		json
// @Produce 	json
// @Param		id path int true "Movie Id"
// @Success		200 {object} data.Movie
// @Router	/v1/movies/{id} [get]
func (app *application) showMovieHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Movie no found")
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   192,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	return c.JSON(http.StatusOK, movie)
}
