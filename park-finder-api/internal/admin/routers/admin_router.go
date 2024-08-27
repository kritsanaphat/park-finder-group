package routers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (ar AdminRouters) adminGetAreaHandler(c echo.Context) error {
	ctx := c.Request().Context()
	status := c.QueryParam("status")
	if status == "waiting" {
		status = "apply completed"
	}

	area := ar.AdminServices.AdminGetParkingArea(ctx, status)

	return c.JSON(http.StatusOK, models.SearchAreaResponse{Data: area, Status: http.StatusOK})
}

func (ar AdminRouters) adminUpdateAreaHandler(c echo.Context) error {
	ctx := c.Request().Context()
	parking_id := c.QueryParam("parking_id")
	status := c.QueryParam("status")
	description := c.QueryParam("description")

	parking := ar.AdminServices.AdminCheckParkingAreaDetail(ctx, parking_id)
	err := ar.AdminServices.AdminUpdateParkingArea(ctx, status, parking_id, description)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	go ar.NotificationService.ParkingAreaStatusUpdateNotification(ctx, parking.ProviderID, status, parking.ParkingName)

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (ar AdminRouters) addRewardHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token naa"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.AddRewardRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	ctx := c.Request().Context()

	_id, err := ar.AdminServices.AddReward(ctx, request, email)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	go ar.NotificationService.AddRewardNotification(ctx, request.Title, request.Description, _id, request.PreviewImageURL)

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (ar AdminRouters) adminGetTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	month := c.QueryParam("month")
	year := c.QueryParam("year")

	data, err := ar.AdminServices.AdminGetTransaction(ctx, year, month)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.AdminTransactionResponseList{Data: data, Status: http.StatusOK})
}

func (ar AdminRouters) adminSubmitReceipt(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(models.SubmitReceiptRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	err := ar.AdminServices.AdminSubmitReceipt(ctx, request.ReceiptImageUrl, request.ProviderID, request.Month, request.Year, request.Price)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})

	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})

}
