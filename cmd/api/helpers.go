package main

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"greenlight.aeropagz.de/internal/validator"
)

func (app *application) readString(c echo.Context, key string, defaultValue string) string {
	s := c.QueryParam(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(c echo.Context, key string, defaultValue []string) []string {
	csv := c.QueryParam(key)
	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(c echo.Context, key string, defaultValue int, v *validator.Validator) int {
	n := c.QueryParam(key)
	if n == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(n)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}

	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error().Msg("error on background task")
			}
		}()

		fn()
	}()
}
