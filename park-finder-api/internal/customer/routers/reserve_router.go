package routers

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (cr CustomerRouters) reserveParkingHandler(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.ReserveRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	order_id := utility.GenerateOrderID()
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is not exists"})
	}

	if account.Fine.Price > 0 {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "You Can not reseve because you stil have fine"})

	}
	if request.Type == "current" {
		order_id, err, id := cr.ReserveService.ReserveParking(ctx, email, order_id, request)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
		}
		return c.JSON(http.StatusOK, models.ReserveResponse{ID: *id, OrderID: order_id, Status: http.StatusOK})

	} else if request.Type == "in_advance" {
		order_id, err, id := cr.ReserveService.ReserveParkingInAdvance(ctx, email, order_id, request)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
		}
		return c.JSON(http.StatusOK, models.ReserveResponse{ID: *id, OrderID: order_id, Status: http.StatusOK})

	}

	return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "invalid type"})
}

func (cr CustomerRouters) startReserveParkingHandler(c echo.Context) error {
	ctx := c.Request().Context()
	parking_id := c.QueryParam("parking_id")
	module_code := c.QueryParam("module_code")
	action := c.QueryParam("action")

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}

	email := utility.GetEmailFromToken(token.Raw)
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Account does not exist", Status: http.StatusUnprocessableEntity})

	}

	reserve := cr.ReserveService.CheckOrderIDByParkingIDAndCustomerEmail(ctx, parking_id, email)
	if reserve == nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Reservation does not exist", Status: http.StatusUnprocessableEntity})

	}

	car := cr.CustomerServices.CheckCustomerCarDetail(ctx, reserve.CarID.Hex())
	if car == nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Car does not exist", Status: http.StatusUnprocessableEntity})

	}
	err, is_found := cr.ReserveService.StartReserveParking(ctx, account.ID, reserve.OrderID, module_code, parking_id, car.LicensePlate, action, reserve.ParkingName)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Order ID and Provider Email Do Not Match, or Payment Was Unsuccessful", Status: http.StatusUnprocessableEntity})

	}
	if is_found {
		return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})

	}
	return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "ok", Status: http.StatusBadRequest})

}

func (cr CustomerRouters) myReserveParkingHandler(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.MyReserveRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	if request.Status == "on_working" || request.Status == "fail" || request.Status == "successful" {
		reservation, err := cr.ReserveService.MyReserve(ctx, email, request.ParkingID, request.Status)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
		}
		if reservation == nil {
			return c.JSON(http.StatusOK, models.MessageResponse{Message: "Not found Reservation log"})
		}

		var Reservation models.Reservations
		Reservation = reservation
		sort.Sort(Reservation)

		var reservation_sort []models.Reservation
		reservation_sort = Reservation

		var mappedReservations []echo.Map
		for _, reservation := range reservation_sort {
			mappedReservations = append(mappedReservations, reservation.ToMapMyReservation())
		}
		return c.JSON(http.StatusOK, models.MyReserveResponse{Reservation: mappedReservations, Status: http.StatusOK})
	}

	return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "invalid type"})
}

func (cr CustomerRouters) myReserveParkingDetailHandler(c echo.Context) error {
	ctx := c.Request().Context()
	reserve_id := c.QueryParam("reserve_id")

	request := new(models.MyReserveRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	reserve := cr.ReserveService.FindReserveDetail(ctx, reserve_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid order id"})

	}
	parking := cr.ReserveService.FindParkingDetail(ctx, reserve.ParkingID.Hex())
	if parking == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid parking_id"})
	}

	response := models.ReserveDetail{
		ID:          reserve.ID,
		OrderId:     reserve.OrderID,
		ProviderID:  reserve.ProviderID,
		ParkingName: parking.ParkingName,
		DateStart:   reserve.DateStart,
		DateEnd:     reserve.DateEnd,
		HourStart:   reserve.HourStart,
		HourEnd:     reserve.HourEnd,
		MinStart:    reserve.MinStart,
		MinEnd:      reserve.MinEnd,
		Latitude:    parking.Address.Latitude,
		Longitude:   parking.Address.Longitude,
	}

	return c.JSON(http.StatusOK, models.ReserveDetailResponse{Data: response, Status: 200})
}

func (cr CustomerRouters) reportParkingVerify(c echo.Context) error {
	ctx := c.Request().Context()
	provider_id := c.QueryParam("provider_id")
	order_id := c.QueryParam("order_id")

	_id, err := primitive.ObjectIDFromHex(provider_id)
	if err != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error convert string to object"})
	}
	reserve := cr.ReserveService.FindReserveDetailByOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "Not found order id"})
	}
	go cr.NotificationService.ReportParkingAreaNotification(ctx, _id, reserve)
	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) reportReserve(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	user := cr.CustomerServices.CheckExistEmail(ctx, email)
	if user == nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found user id", Status: http.StatusBadRequest})
	}

	request := new(models.ReportRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	reserve := cr.ReserveService.CheckExistOrderID(ctx, request.OrderID)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found order id", Status: http.StatusBadRequest})
	}

	err := cr.ReserveService.ReportReservation(ctx, user.ID, reserve.ProviderID, request.Content, request.OrderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})

	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})

}

func (cr CustomerRouters) createReviewHandler(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.ReviewRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	user := cr.CustomerServices.CheckExistEmail(ctx, email)

	err := cr.ReserveService.CreateReview(ctx, user.IDToString(), user.FirstName, user.LastName, request.Comment, request.ParkingID, request.OrderID, request.ReviewScore)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	err = cr.CustomerServices.RemoveReviewCache(request.OrderID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) checkCanReview(c echo.Context) error {
	ctx := context.Background()
	order_id := c.QueryParam("order_id")

	can_review := cr.CustomerServices.CheckExistReview(order_id)
	if can_review {
		return c.JSON(http.StatusOK, models.StatusResponseBool{Message: true, Status: http.StatusOK})

	} else {
		reserve := cr.ReserveService.FindReserveDetailByOrderID(ctx, order_id)
		if reserve == nil {
			return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found order id", Status: http.StatusBadRequest})
		}
		parking_area := cr.ReserveService.FindParkingDetail(ctx, reserve.ParkingID.Hex())
		if parking_area == nil {
			return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found parking area", Status: http.StatusBadRequest})
		}

		review_list := parking_area.Review
		if review_list == nil {
			return c.JSON(http.StatusOK, models.StatusResponseBool{Message: false, Status: http.StatusOK})
		}
		my_review := new(models.ReviewParkingArea)
		for _, temp := range review_list {
			if temp.OrderID == order_id {
				my_review = &temp
				break
			} else {
				my_review = nil
			}
		}
		if my_review == nil {
			return c.JSON(http.StatusOK, models.StatusResponseBool{Message: false, Status: http.StatusOK})
		}
		return c.JSON(http.StatusOK, models.StatusReviewResponse{Message: *my_review, Status: http.StatusOK})

	}
}

func (cr CustomerRouters) checkFine(c echo.Context) error {

	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	user := cr.CustomerServices.CheckExistEmail(ctx, email)
	if user == nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found user id", Status: http.StatusBadRequest})
	}
	if user.Fine.Price > 0 {
		return c.JSON(http.StatusOK, models.StatusResponseFine{Data: user.Fine, Status: http.StatusOK})
	} else {
		return c.JSON(http.StatusOK, models.StatusResponseFine{Data: nil, Status: http.StatusOK})

	}
}

func (cr CustomerRouters) extendReserveHandler(c echo.Context) error {
	ctx := c.Request().Context()
	order_id := c.QueryParam("order_id")
	action := c.QueryParam("action")

	reserve := cr.ReserveService.CheckExistOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found order id", Status: http.StatusBadRequest})
	}

	account := cr.CustomerServices.CheckExistEmail(ctx, reserve.CustomerEmail)
	if account == nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not account order id", Status: http.StatusBadRequest})
	}

	reserves_dulplicate, err := cr.ReserveService.CheckReserveInNextHour(ctx, reserve.ParkingID, reserve.HourEnd, reserve.HourEnd+1, reserve.DateEnd)
	if err != nil {
		return err
	}
	if action == "automatic" {
		var price = reserve.Price
		if len(*reserves_dulplicate) > 0 {
			price = reserve.Price * 2
			fmt.Println("is have reserves_dulplicate ")
			err = cr.ReserveService.ExtendReserve(ctx, order_id, action, reserve.DateEnd, reserve.HourEnd+1, reserve.MinEnd)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
			}
			for _, temp := range *reserves_dulplicate {
				err = cr.CustomerServices.CustomerRefund(ctx, temp.CustomerEmail, temp.Price)
				if err != nil {
					return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
				}
				err = cr.ReserveService.UpdateReserveStatusAndRemoveJob(ctx, temp.OrderID)
				if err != nil {
					return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
				}
				go cr.NotificationService.ReservationCancelNotification(ctx, account.ID, &temp)
			}
		} else {
			fmt.Println("is not have reserves_dulplicate ")
			err = cr.ReserveService.ExtendReserve(ctx, order_id, action, reserve.DateEnd, reserve.HourEnd+1, reserve.MinEnd)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
			}
		}

		err := cr.CustomerServices.CustomerFine(ctx, account.ID, price, reserve)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
		}
		return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})

	} else if action == "normal" {
		var price = reserve.Price
		if len(*reserves_dulplicate) > 0 {
			return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Cann't extend becasue have a customer reserve in a next hour", Status: http.StatusBadRequest})

		} else {
			//ไม่มีผู้ได้รับผลกระทบ
			fmt.Println("is not have reserves_dulplicate ")
			err = cr.ReserveService.ExtendReserve(ctx, order_id, action, reserve.DateEnd, reserve.HourEnd+1, reserve.MinEnd)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
			}
		}
		var cashback = 0
		if account.Cashback >= price {
			cashback = price
		}
		return_to_line_reserve := models.LineReserveRequest{
			OrderID:     reserve.OrderID,
			ProviderID:  reserve.ParkingID.Hex(),
			ParkingID:   reserve.ParkingID.Hex(),
			Quantity:    1,
			Price:       reserve.Price,
			ParkingName: reserve.ParkingName,
			CashBack:    cashback,
		}
		return c.JSON(http.StatusOK, models.ExtendReserveResponse{Data: return_to_line_reserve, Status: http.StatusOK})
	}
	return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "invalid type"})

}

func (cr CustomerRouters) getParkingDetailHandler(c echo.Context) error {
	ctx := c.Request().Context()

	_, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}

	parking_id := c.QueryParam("parking_id")

	parking := cr.ReserveService.FindParkingDetail(ctx, parking_id)
	if parking != nil {
		return c.JSON(http.StatusOK, models.ParkingDetailResponse{Data: parking, Status: http.StatusOK})
	}

	return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found Parking id"})
}

func (cr CustomerRouters) captureCarReserve(c echo.Context) error {
	ctx := c.Request().Context()
	order_id := c.QueryParam("order_id")
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found account id", Status: http.StatusBadRequest})
	}
	reserve := cr.ReserveService.CheckExistOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found order id", Status: http.StatusBadRequest})
	}

	if reserve.ModuleCode == "" {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: "Not found module code", Status: http.StatusBadRequest})
	}

	imge_url, err := cr.ReserveService.CaptureCarReserve(ctx, account.ID, reserve.ModuleCode, reserve.ParkingName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: imge_url, Status: http.StatusOK})
}
