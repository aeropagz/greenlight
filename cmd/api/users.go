package main

import (
	"errors"
	"net/http"
	"time"

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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.html", data)
		if err != nil {
			app.logger.Error().AnErr("email", err).Msg("could not send register mail")
		}
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// @Summary     Activate an a user
// @Tags		User
// @Accept 		json
// @Produce 	json
// @Param		request body main.activateUserHandler.Input true "Activation Token"
// @Success		200 {object} data.User
// @Router	/v1/users/activated [put]
func (app *application) activateUserHandler(c echo.Context) error {
	type Input struct {
		TokenPlaintext string `json:"token"`
	}

	input := &Input{}
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, v.Errors)
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired token")
			return echo.NewHTTPError(http.StatusUnprocessableEntity, v.Errors)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	user.Activated = true

	user, err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)

}
