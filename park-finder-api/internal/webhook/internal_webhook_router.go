package routers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
)

// Cronjob  CallBack Method
func (ws WebhookRouters) BeforeTimeOutReserve(c echo.Context) error {
	ctx := c.Request().Context()
	order_id := c.QueryParam("order_id")
	reserve := ws.ReserveServices.CheckExistOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found User ID"})
	}
	customer := ws.CustomerServices.CheckExistEmail(ctx, reserve.CustomerEmail)
	car := ws.CustomerServices.CheckCustomerCarDetail(ctx, reserve.CarID.Hex())
	if car == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found Car ID"})
	}
	is_found := ws.HttpClient.SendCameraServiceToDectectIncommingCar(reserve.ModuleCode, car.LicensePlate)
	fmt.Println("status detect is ", is_found)

	if !is_found {
		err := ws.ReserveServices.UpdateReserveStatus(ctx, "Successful", order_id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Error form update Successful Status reservation"})

		}
		err = ws.CacheReview(order_id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})

		}
		go ws.HttpClient.SendRemoveJobAPI("ATOR_" + order_id)
		go ws.HttpClient.SendRemoveJobAPI("TOR_" + order_id)
		go ws.NotificationService.LeaveTimeOutReserveNotification(ctx, customer.ID, reserve.ParkingID, order_id, reserve.ParkingName)

	} else {
		reserve := ws.ReserveServices.FindReserveDetailByOrderID(ctx, order_id)
		if reserve == nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not found order id"})
		}
		go ws.NotificationService.BeforeTimeOutReserveNotification(ctx, customer.ID, reserve)
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "ok"})
}

func (ws WebhookRouters) TimeOutReserve(c echo.Context) error {
	ctx := c.Request().Context()
	order_id := c.QueryParam("order_id")
	reserve := ws.ReserveServices.CheckExistOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found User ID"})
	}
	customer := ws.CustomerServices.CheckExistEmail(ctx, reserve.CustomerEmail)
	if customer == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found customer ID"})
	}
	car := ws.CustomerServices.CheckCustomerCarDetail(ctx, reserve.CarID.Hex())
	if car == nil {
		fmt.Println(car)
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found Car ID"})
	}
	is_found := ws.HttpClient.SendCameraServiceToDectectIncommingCar(reserve.ModuleCode, car.LicensePlate)
	fmt.Println("status detect is ", is_found)

	if !is_found {
		err := ws.ReserveServices.UpdateReserveStatus(ctx, "Successful", order_id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Error form update Successful Status reservation"})

		}
		err = ws.CacheReview(order_id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})

		}
		go ws.HttpClient.SendRemoveJobAPI("ATOR_" + order_id)
		go ws.NotificationService.LeaveTimeOutReserveNotification(ctx, customer.ID, reserve.ParkingID, order_id, reserve.ParkingName)

	} else {
		reserve := ws.ReserveServices.FindReserveDetailByOrderID(ctx, order_id)
		if reserve == nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not found order id"})
		}
		ws.NotificationService.TimeOutReserveNotification(ctx, customer.ID, reserve)
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "ok"})
}

func (ws WebhookRouters) AfterTimeOutReserve(c echo.Context) error {
	ctx := c.Request().Context()
	order_id := c.QueryParam("order_id")
	reserve := ws.ReserveServices.CheckExistOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found User ID"})
	}
	customer := ws.CustomerServices.CheckExistEmail(ctx, reserve.CustomerEmail)
	if customer == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found customer ID"})
	}
	car := ws.CustomerServices.CheckCustomerCarDetail(ctx, reserve.CarID.Hex())
	if car == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found Car ID"})
	}
	is_found := ws.HttpClient.SendCameraServiceToDectectIncommingCar(reserve.ModuleCode, car.LicensePlate)
	fmt.Println("status detect is ", is_found)

	if !is_found {
		err := ws.ReserveServices.UpdateReserveStatus(ctx, "Successful", order_id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Error form update Successful Status reservation"})

		}
		err = ws.CacheReview(order_id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})

		}
		go ws.NotificationService.LeaveTimeOutReserveNotification(ctx, customer.ID, reserve.ParkingID, order_id, reserve.ParkingName)

	} else {
		is_not_err := ws.HttpClient.SendToAutomaticExtendReserve(order_id)
		if is_not_err {
			customer := ws.CustomerServices.CheckExistEmail(ctx, reserve.CustomerEmail)
			if customer == nil {
				return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not Found customer ID"})
			}
			ws.NotificationService.AfterTimeOutReserveNotification(ctx, customer.ID, reserve, customer.Fine.Price)
		}

	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "ok"})
}

// Notification CallBack Method
func (ws WebhookRouters) ConfirmReserveInAdvance(c echo.Context) error {
	ctx := c.Request().Context()
	order_id := c.QueryParam("order_id")
	email := c.QueryParam("email")

	account := ws.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not found customer email"})
	}
	err := ws.ReserveServices.UpdateReserveStatus(ctx, "Process", order_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	reserve := ws.ReserveServices.FindReserveDetailByOrderID(ctx, order_id)
	if reserve == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not found order id"})
	}
	parking_area := ws.ReserveServices.FindParkingDetail(ctx, reserve.ParkingID.Hex())
	if parking_area == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not found parking_area"})
	}

	full_address := parking_area.Address.AddressText + parking_area.Address.Sub_district + parking_area.Address.District + parking_area.Address.Province + parking_area.Address.Postal_code

	go ws.NotificationService.ProviderConfirmReserveInAdvanceNotification(ctx, account.ID, reserve, full_address)
	go ws.HttpClient.SendRemoveJobAPI("CRIA_" + order_id + "," + email)

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "ok"})
}

func (ws WebhookRouters) CancelReserveInAdvance(c echo.Context) error {
	ctx := c.Request().Context()
	order_id := c.QueryParam("order_id")
	email := c.QueryParam("email")

	account := ws.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not found customer email"})
	}

	err := ws.ReserveServices.UpdateReserveStatus(ctx, "Cancel", order_id)
	if err != nil {
		return err
	}

	err = ws.refundCashBack(ctx, order_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})

	}
	reserve := ws.ReserveServices.FindReserveDetailByOrderID(ctx, order_id)
	if reserve == nil {
		return err
	}
	parking_area := ws.ReserveServices.FindParkingDetail(ctx, reserve.ParkingID.Hex())

	if parking_area == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Not found parking_area"})
	}

	full_address := parking_area.Address.AddressText + parking_area.Address.Sub_district + parking_area.Address.District + parking_area.Address.Province + parking_area.Address.Postal_code

	go ws.NotificationService.ProviderCancelReserveInAdvanceNotification(ctx, account.ID, reserve, full_address)

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "ok"})
}

func (ws WebhookRouters) refundCashBack(ctx context.Context, order_id string) error {
	reserve := ws.ReserveServices.CheckExistOrderID(ctx, order_id)
	if reserve == nil {
		return errors.New("not found reserve id")
	}
	err := ws.CustomerServices.CustomerRefund(ctx, reserve.CustomerEmail, reserve.Price)
	if err != nil {
		return err
	}

	return nil
}

func (ws WebhookRouters) CacheReview(order_id string) error {
	key := order_id

	// Cache for 1 hour
	result := ws.Redis.Set(key, true, 3600)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user review for order id:", order_id)

	return nil
}
