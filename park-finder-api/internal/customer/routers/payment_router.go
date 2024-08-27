package routers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (cr CustomerRouters) linePayReserveHandler(c echo.Context) error {

	ctx := c.Request().Context()
	action := c.QueryParam("action")

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	request := new(models.LineReserveRequest)

	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "this account doesn't exist"})
	}

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	reserve := cr.ReserveService.CheckExistOrderID(ctx, request.OrderID)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "this order id doesn't exist"})
	}
	//จ่าย cashback เต็มจำนวน
	if request.CashBack == request.Price*int(request.Quantity) {
		if account.Cashback >= request.CashBack {
			err := cr.PaymentService.CashbackReserve(ctx, email, request.OrderID, request.ParkingID, reserve.Type, action, request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
			}
			err = cr.CustomerServices.CustomerUpdateCashback(ctx, email, account.Cashback-request.CashBack)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
			}

			return c.JSON(http.StatusOK, models.StatusResponse{Message: "Success payment with cashback", Status: http.StatusOK})
		} else {
			return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Your cashback not enough", Status: http.StatusBadRequest})
		}

	} else {
		var webUrl string
		if action != "" {
			var err error = nil
			fmt.Println("Fine")
			webUrl, err = cr.PaymentService.LinePayReserve(ctx, email, request.OrderID, request.ParkingID, "fine", request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
			}
		} else {
			fmt.Println("Not fine")
			var err error = nil
			webUrl, err = cr.PaymentService.LinePayReserve(ctx, email, request.OrderID, request.ParkingID, "", request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
			}
		}

		if request.CashBack > 0 {
			err := cr.CustomerServices.CustomerUpdateCashback(ctx, email, account.Cashback-request.CashBack)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
			}
		}
		return c.JSON(http.StatusOK, models.StatusResponsePayment{Message: webUrl, ReservationID: reserve.ID.Hex(), Status: http.StatusOK})
	}

}
