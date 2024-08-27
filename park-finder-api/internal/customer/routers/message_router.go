package routers

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (cr CustomerRouters) customerChatLogHandler(c echo.Context) error {

	ctx := c.Request().Context()
	reservation_id := c.QueryParam("reservation_id")
	start := c.QueryParam("start")
	limit := c.QueryParam("limit")

	start_int, err := strconv.Atoi(start)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid start parameter")
	}

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid start parameter")
	}

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid email"})

	}
	chat_list := cr.MessageServices.RetriveChatLog(ctx, reservation_id, start_int, limit_int)
	if chat_list == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Error form retrive chat log"})

	}

	return c.JSON(http.StatusOK, models.ListMessageLogResponse{Status: 200, Data: chat_list})
}
