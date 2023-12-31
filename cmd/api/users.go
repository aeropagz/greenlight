package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"greenlight.aeropagz.de/internal/data"
	"greenlight.aeropagz.de/internal/validator"
)

// @Summary     Create an a user
// @Tags		User
// @Accept 		json
// @Produce 	json
// @Param		request body main.createUserHandler.userIn true "Movie creation dto"
// @Success		201 {object} data.User
// @Router	/v1/users [post]
func (app *application) createUserHandler(c echo.Context) error {
	type userIn struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := &userIn{}

	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, v.Errors)
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			return echo.NewHTTPError(http.StatusUnprocessableEntity, v.Errors)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	err = app.mailer.Send(user.Email, "user_welcome.html", user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}


	return c.JSON(http.StatusCreated, user)
}
