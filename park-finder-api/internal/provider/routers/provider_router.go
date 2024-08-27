package routers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (pr ProviderRouters) providerRegisterAreaLocationHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.RegisterParkingAreaFirstStepRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	context := c.Request().Context()

	verify := pr.ProviderServices.CheckVerifyEmail(context, email)
	if verify == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is not verify"})
	}

	_id, err := pr.ProviderServices.ProviderRegisterAreaLocaion(context, request, email)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.RegisterAreaLocationResponse{ParkingID: _id, Message: "ok", Status: http.StatusOK})
}

func (pr ProviderRouters) providerRegisterAreaDocumentHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	parking_id := c.QueryParam("parking_id")

	request := new(models.RegisterParkingAreaDocumentStepRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	context := c.Request().Context()

	verify := pr.ProviderServices.CheckVerifyEmail(context, email)
	if verify == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is not verify"})
	}

	if err := pr.ProviderServices.ProviderRegisterAreaDocument(context, request, email, parking_id); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (pr ProviderRouters) updateProviderProfileHandler(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.UpdateProfileRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	if err := pr.ProviderServices.UpdateProviderProfile(ctx, email, request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (pr ProviderRouters) updateBankAccountHandler(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.BankAccount)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	if err := pr.ProviderServices.UpdateProviderBankAccount(ctx, email, request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (pr ProviderRouters) getProviderProfileHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	ctx := c.Request().Context()
	account := pr.ProviderServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile emails"})
	}
	return c.JSON(http.StatusOK, models.CustomerProfileResponse{Profile: account.ToMapProfile(), Status: http.StatusOK})
}

func (pr ProviderRouters) getProviderAreaHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	context := c.Request().Context()
	account := pr.ProviderServices.CheckExistEmail(context, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile token"})
	}

	areas := pr.ProviderServices.GetProviderArea(context, account.IDToString())
	if areas == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get areas by _id"})
	}
	sortAreasByStatusApply(areas)

	return c.JSON(http.StatusOK, models.SearchAreaResponse{Data: &areas, Status: http.StatusOK})
}

func (pr ProviderRouters) getParkingDetailHandler(c echo.Context) error {
	ctx := c.Request().Context()

	_, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}

	parking_id := c.QueryParam("parking_id")

	parking := pr.ReserveService.FindParkingDetail(ctx, parking_id)
	if parking != nil {
		return c.JSON(http.StatusOK, models.ParkingDetailResponse{Data: parking, Status: http.StatusOK})
	}

	return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "invalid type"})
}

func (pr ProviderRouters) getProviderProfitHandler(c echo.Context) error {

	ctx := c.Request().Context()
	query_type := c.QueryParam("query_type")
	parking_id := c.QueryParam("list_parking_id")

	list_parking_id := strings.Split(parking_id, ",")

	if query_type == "daily" {
		response, count, err := pr.ProviderServices.GetProviderProfitDaily(ctx, list_parking_id)
		if err != nil {
			return c.JSON(http.StatusConflict, models.MessageResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, models.WeeklyProfitResponse{Data: []models.DailyProfitResponse{{Date: response.Date, Sum: response.Sum}}, Count: count})
	}
	if query_type == "weekly" {
		sum, count, err := pr.ProviderServices.GetProviderProfitWeekly(ctx, list_parking_id)
		if err != nil {
			return c.JSON(http.StatusConflict, models.MessageResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, models.WeeklyProfitResponse{Data: *sum, Count: count})
	}
	if query_type == "monthly" {
		sum, count, err := pr.ProviderServices.GetProviderProfitMontly(ctx, list_parking_id)
		if err != nil {
			return c.JSON(http.StatusConflict, models.MessageResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, models.WeeklyProfitResponse{Data: *sum, Count: count})
	}
	if query_type == "yearly" {
		sum, count, err := pr.ProviderServices.GetProviderProfitYearly(ctx, list_parking_id)
		if err != nil {
			return c.JSON(http.StatusConflict, models.MessageResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, models.WeeklyProfitResponse{Data: *sum, Count: count})
	}

	return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid query_type"})
}

func (pr ProviderRouters) updateOpenAreaQuickStatusHandler(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := pr.ProviderServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile token"})
	}
	if account.BankAccount.AccountBookImageUrl == "" {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "bank information is null"})
	}

	request := new(models.UpdateOpenAreaQuickStatusRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	provider_area := pr.ProviderServices.CheckValidProviderArea(ctx, account.IDToString(), request.ParkingAreaID)
	if provider_area != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "parking area & provider not match"})
	}

	if request.Type == "normal" {
		if request.Status == "close" {
			if reservation := pr.ProviderServices.CheckQuickAvailabilityCloseProviderArea(ctx, request.ParkingAreaID, request.Range); reservation != nil {
				return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Cann't close area because there were customers during that time.", Status: http.StatusUnprocessableEntity})
			}
			if err := pr.ProviderServices.UpdateOpenAreaQuickStatus(ctx, request.ParkingAreaID, request.Status, request.Range); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
			}
		} else {
			if err := pr.ProviderServices.UpdateOpenAreaQuickStatus(ctx, request.ParkingAreaID, request.Status, request.Range); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
			}
		}
	} else if request.Type == "force" {
		if request.Status == "close" {
			if reservation := pr.ProviderServices.CheckQuickAvailabilityCloseProviderArea(ctx, request.ParkingAreaID, request.Range); reservation != nil {
				for _, temp := range *reservation {
					err := pr.CustomerServices.CustomerRefund(ctx, temp.CustomerEmail, temp.Price+int((float32(temp.Price)*0.1)))
					if err != nil {
						return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
					}
					err = pr.ProviderServices.FineProvider(ctx, temp.ProviderID, temp.Price-int(float32(temp.Price)*0.65))
					if err != nil {
						return err
					}
					err = pr.ReserveService.UpdateReserveStatus(ctx, "Cancel", temp.OrderID)
					if err != nil {
						return err
					}
					customer_account := pr.CustomerServices.CheckExistEmail(ctx, temp.CustomerEmail)
					if customer_account == nil {
						return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Not found customer account id", Status: http.StatusUnprocessableEntity})

					}
					go pr.HttpClient.SendRemoveJobAPI("BTOR_" + temp.OrderID)
					go pr.HttpClient.SendRemoveJobAPI("TOR_" + temp.OrderID)
					go pr.HttpClient.SendRemoveJobAPI("ATOR_" + temp.OrderID)
					go pr.NotificationService.ReservationCancelNotification(ctx, customer_account.ID, &temp)
				}
			}
			if err := pr.ProviderServices.UpdateOpenAreaQuickStatus(ctx, request.ParkingAreaID, request.Status, request.Range); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
			}
		} else {
			if err := pr.ProviderServices.UpdateOpenAreaQuickStatus(ctx, request.ParkingAreaID, request.Status, request.Range); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
			}
		}
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (pr ProviderRouters) updateOpenAreaDailyStatusHandler(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := pr.ProviderServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile token"})
	}
	if account.BankAccount.AccountBookImageUrl == "" {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "bank information is null"})
	}

	request := new(models.UpdateOpenAreaDailyStatusRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	provider_area := pr.ProviderServices.CheckValidProviderArea(ctx, account.IDToString(), request.ParkingAreaID)
	if provider_area != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "parking area & provider not match"})
	}

	if request.Type == "normal" {

		if reservation := pr.ProviderServices.CheckDailyAvailabilityCloseProviderArea(ctx, *request); reservation != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Cann't close area because changing daily open status must not have any active customers in change time", Status: http.StatusUnprocessableEntity})
		}

		if err := pr.ProviderServices.UpdateOpenAreaDailyStatus(ctx, request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
		}
	}

	if request.Type == "force" {
		if reservation := pr.ProviderServices.CheckDailyAvailabilityCloseProviderArea(ctx, *request); reservation != nil {
			for _, temp := range *reservation {
				err := pr.CustomerServices.CustomerRefund(ctx, temp.CustomerEmail, temp.Price+int((float32(temp.Price)*0.1)))
				if err != nil {
					return c.JSON(http.StatusBadRequest, models.StatusResponse{Message: err.Error(), Status: http.StatusBadRequest})
				}

				err = pr.ProviderServices.FineProvider(ctx, temp.ProviderID, temp.Price-int(float32(temp.Price)*0.65))
				if err != nil {
					return err
				}
				err = pr.ReserveService.UpdateReserveStatus(ctx, "Cancel", temp.OrderID)
				if err != nil {
					return err
				}
				customer_account := pr.CustomerServices.CheckExistEmail(ctx, temp.CustomerEmail)
				if customer_account == nil {
					return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: "Not found customer account id", Status: http.StatusUnprocessableEntity})

				}

				go pr.HttpClient.SendRemoveJobAPI("BTOR_" + temp.OrderID)
				go pr.HttpClient.SendRemoveJobAPI("TOR_" + temp.OrderID)
				go pr.HttpClient.SendRemoveJobAPI("ATOR_" + temp.OrderID)
				go pr.NotificationService.ReservationCancelNotification(ctx, customer_account.ID, &temp)
			}

		}
		if err := pr.ProviderServices.UpdateOpenAreaDailyStatus(ctx, request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})

		}
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (pr ProviderRouters) updateParkingAreaPriceHandler(c echo.Context) error {
	ctx := c.Request().Context()
	parking_id := c.QueryParam("parking_id")
	price := c.QueryParam("price")
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := pr.ProviderServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile token"})
	}
	price_int, err := strconv.Atoi(price)
	if err != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "invalid price, price must be number"})
	}
	err = pr.ProviderServices.CheckUpdatePriceArea(ctx, parking_id, int16(price_int))
	if err != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func sortAreasByStatusApply(areas []models.ParkingArea) {
	sort.Slice(areas, func(i, j int) bool {
		// ถ้า StatusApply เป็น "accepted" ให้ข้อมูลนี้ไปขึ้นก่อน
		if areas[i].StatusApply == "accepted" && areas[j].StatusApply != "accepted" {
			return true
		}
		return false
	})
}
