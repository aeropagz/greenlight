package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"greenlight.aeropagz.de/internal/data"
	"greenlight.aeropagz.de/internal/validator"
)

// @Summary     Create an a movie
// @Tags		Movie
// @Accept 		json
// @Produce 	json
// @Param		request body main.createMovieHandler.input true "Movie creation dto"
// @Success		200 {object} data.Movie
// @Router	/v1/movies [post]
func (app *application) createMovieHandler(c echo.Context) error {
	type movieCreation struct {
		Title   string   `json:"title"`
		Year    int32   `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	input := new(movieCreation)

	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	movie := &data.Movie{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, v.Errors)
	}

	return c.JSON(http.StatusOK, input)
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
