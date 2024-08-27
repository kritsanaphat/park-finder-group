package routers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (cr CustomerRouters) customerNotificatiobListHandler(c echo.Context) error {

	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid email"})
	}

	notification_list, err := cr.NotificationService.NotificationList(ctx, account.ID.Hex(), "customer")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Error form retrive chat list"})
	}

	return c.JSON(http.StatusOK, models.ListNotificationResponse{Status: 200, Data: &notification_list})
}
