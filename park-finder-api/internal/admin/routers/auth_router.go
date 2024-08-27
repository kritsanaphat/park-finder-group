package routers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (ar AdminRouters) adminRegisterHandler(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(models.RegisterAccountRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := ar.AdminServices.CheckExistEmail(ctx, request.Email)
	if account != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is already exists"})
	}

	if err := ar.AdminServices.AdminRegister(ctx, request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (ar AdminRouters) adminLoginHandler(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(models.LoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := ar.AdminServices.CheckExistEmail(ctx, request.Email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is not registered"})
	}

	if !ar.AdminServices.CheckPassword(*account, request.Password) {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Incorrect password"})
	}

	accessToken := utility.GenerateToken(request.Email, "admin")
	if accessToken == "" {
		return c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Cannot generate new token"})
	}

	err := ar.AdminServices.AddToken(ctx, account, accessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Error while add user token"})
	}

	return c.JSON(http.StatusOK, models.LoginResponse{AccessToken: accessToken, Data: account.ToMap()})
}

func (cr AdminRouters) adminLogoutHandler(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}

	email := utility.GetEmailFromToken(token.Raw)
	user := cr.AdminServices.CheckExistEmail(ctx, email)
	if user == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "User not found"})
	}

	err := cr.AdminServices.RevokeToken(ctx, user, token.Raw)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Error while revoke user token"})
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "Logout succuess"})
}

// errorHandler handle expired jwt token by update token status and stamp revoke date
func (ar AdminRouters) errorHandler(c echo.Context, err error) error {
	message := "invalid or expired jwt"
	switch err.(type) {
	case *echojwt.TokenExtractionError:
		message = "missing or malformed jwt"
	default:

		ctx := c.Request().Context()
		auth := c.Request().Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			message = "missing or malformed jwt"
			return c.JSON(http.StatusUnauthorized, models.MessageResponse{
				Message: message,
			})
		}
		tokenStr := strings.Split(auth, " ")[1]
		token := ar.AdminServices.CheckExistToken(ctx, tokenStr)
		if token != nil && token.Valid {
			expire := utility.GetExpireDateFromToken(tokenStr)
			expireDate := time.Unix(expire, 0)
			ar.AdminServices.RevokeExpireToken(ctx, tokenStr, expireDate)
		}
	}
	return c.JSON(http.StatusUnauthorized, models.MessageResponse{
		Message: message,
	})
}

// parseToken validate the given auth token, Return token if recieved token is valid otherwise return nil
func (ar AdminRouters) parseToken(c echo.Context, auth string) (interface{}, error) {
	token, err := jwt.Parse(auth, utility.KeyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	ctx := c.Request().Context()
	dbToken := ar.AdminServices.CheckExistToken(ctx, token.Raw)
	if !dbToken.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	role := utility.GetRoleFromToken(token.Raw)
	if role != "admin" {
		return nil, fmt.Errorf("token invalid role")
	}

	return token, nil
}
