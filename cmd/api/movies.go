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
// @Success		201 {object} data.Movie
// @Router	/v1/movies [post]
func (app *application) createMovieHandler(c echo.Context) error {
	input := new(data.MovieCreation)

	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
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
	if err != nil {
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
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return echo.NewHTTPError(http.StatusConflict, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, movie)
}

// @Summary     Delete an a movie
// @Tags		Movie
// @Accept 		json
// @Produce 	json
// @Param		id path int true "Movie Id"
// @Success		200 {object} data.Message
// @Router	/v1/movies/{id} [delete]
func (app *application) deleteMovieHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Record not found")
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, &data.Message{Message: "movie successfully deleted"})
}

func (app *application) listMovieHandler(c echo.Context) error {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	input.Title = app.readString(c, "title", "")
	input.Genres = app.readCSV(c, "genres", []string{})
	input.Filters.Page = app.readInt(c, "page", 1, v)
	input.Filters.PageSize = app.readInt(c, "page_size", 20, v)
	input.Filters.Sort = app.readString(c, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		return echo.NewHTTPError(http.StatusBadRequest, v.Errors)
	}

	movies, meta, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	envelope := map[string]any{
		"movies": movies,
		"meta":   meta,
	}

	return c.JSON(http.StatusOK, envelope)
}
