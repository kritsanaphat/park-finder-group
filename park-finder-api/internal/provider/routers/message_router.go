package routers

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

// message service
func (pr ProviderRouters) providerChatLogHandler(c echo.Context) error {

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
	account := pr.ProviderServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid email"})

	}
	chat_list := pr.MessageServices.RetriveChatLog(ctx, reservation_id, start_int, limit_int)
	if chat_list == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Error form retrive chat log"})

	}

	return c.JSON(http.StatusOK, models.ListMessageLogResponse{Status: 200, Data: chat_list})
}

func (cr ProviderRouters) providerChatListHandler(c echo.Context) error {

	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := cr.ProviderServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid email"})

	}
	chat_list, err := cr.MessageServices.RetriveChatList(ctx, account.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Error form retrive chat list"})

	}
	msg_room_map := make([]map[string]interface{}, len(chat_list))
	for i, room := range chat_list {
		msg_room_map[i] = room.ToMapMessageRoom()
	}

	return c.JSON(http.StatusOK, models.ListMessageRoomResponse{Status: 200, Data: &msg_room_map})
}
