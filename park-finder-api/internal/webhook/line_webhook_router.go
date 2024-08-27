package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
)

func (ws WebhookRouters) LinePayConfirmCallBackWebhookHandler(c echo.Context) error {
	ctx := c.Request().Context()
	transactionId := c.QueryParam("transactionId")
	order_id := c.QueryParam("orderId")

	reserve := ws.ReserveServices.CheckExistOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "This order id not found"})
	}

	err := ws.PaymentServices.LinePayConfirm(ctx, transactionId, order_id, reserve.Type, reserve.IsExtend)
	if err != nil {
		return err
	}
	if reserve.Type == "in_advance" {
		ws.NotificationService.ConfirmReserveInAdvanceNotification(ctx, reserve.ProviderID, reserve.CustomerEmail, *reserve)
	}

	return nil
}
func (ws WebhookRouters) LinePayConfirmFineCallBackWebhookHandler(c echo.Context) error {
	ctx := c.Request().Context()
	transactionId := c.QueryParam("transactionId")
	order_id := c.QueryParam("orderId")

	err := ws.PaymentServices.LinePayConfirmFine(ctx, transactionId, order_id)
	if err != nil {
		return err
	}

	return nil
}

func (ws WebhookRouters) LinePayCancelCallBackWebhookHandler(c echo.Context) error {
	ctx := c.Request().Context()
	transactionId := c.QueryParam("transactionId")
	order_id := c.QueryParam("orderId")

	err := ws.PaymentServices.LinePayCancel(ctx, transactionId, order_id)
	if err != nil {
		return err
	}
	return nil
}
