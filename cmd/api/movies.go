package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"greenlight.aeropagz.de/internal/data"
	"greenlight.aeropagz.de/internal/validator"
)

// @Summary     Create an a movie
// @Tags		Movie
// @Accept 		json
// @Produce 	json
// @Param		request body data.MovieCreation true "Movie creation dto"
// @Success		200 {object} data.Movie
// @Router	/v1/movies [post]
func (app *application) createMovieHandler(c echo.Context) error {
	input := new(data.MovieCreation)

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

	err := app.models.Movies.Insert(movie)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/v1/movies/%d", movie.ID))
	return c.JSON(http.StatusCreated, input)
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
		return echo.NewHTTPError(http.StatusNotFound, "No record found")
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, movie)
}

// @Summary     Show an a movie
// @Tags		Movie
// @Accept 		json
// @Produce 	json
// @Param		id path int true "Movie Id"
// @Param		request body data.MovieCreation true "Movie creation dto"
// @Success		200 {object} data.Movie
// @Router	/v1/movies/{id} [put]
func (app *application) updateMovieHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Record not found")
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil{
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "unexpected error")
		}
	}

	var input = new(data.MovieCreation)

	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()
	
	if data.ValidateMovie(v, movie); !v.Valid() {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, v.Errors)
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, movie)
}