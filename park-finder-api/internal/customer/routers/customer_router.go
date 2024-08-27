package routers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (cr CustomerRouters) updateCustomerProfileHandler(c echo.Context) error {
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
	if err := cr.CustomerServices.UpdateCustomerProfile(ctx, email, request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) customerProfileHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	ctx := c.Request().Context()
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile email"})
	}
	return c.JSON(http.StatusOK, models.CustomerProfileResponse{Profile: account.ToMapProfile(), Status: http.StatusOK})
}

func (cr CustomerRouters) customerRegisteCarHandler(c echo.Context) error {
	request := new(models.RegisterCarRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	car := cr.CustomerServices.CheckExistCarName(ctx, email, request.Name)
	if car != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This car name is already exists"})
	}

	car = cr.CustomerServices.CheckExistCarLicensePlate(ctx, email, request.LicensePlate)
	if car != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This car license plate is already exists"})
	}

	if err := cr.CustomerServices.CustomerRegisterCar(ctx, request, email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) customerCarlHandler(c echo.Context) error {
	context := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	cars := cr.CustomerServices.CustomerCar(context, email)
	if cars == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "Not Found Car"})
	}

	return c.JSON(http.StatusOK, models.CustomerCarResponse{Data: cars, Status: http.StatusOK})
}

func (cr CustomerRouters) customerAddressHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	ctx := c.Request().Context()
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile email"})
	}
	return c.JSON(http.StatusOK, models.CustomerAddressResponse{Data: account.Address, Status: http.StatusOK})
}

func (cr CustomerRouters) customerRegisterAddressHandler(c echo.Context) error {
	request := new(models.RegisterAddressRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	location_name := cr.CustomerServices.CheckExistLocationName(ctx, email, request.LocationName)

	if location_name != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This location name is already exists"})
	}

	if len(strings.Split(request.Address, " ")) < 6 {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This address is invalid"})
	}

	if err := cr.CustomerServices.CustomerRegisterAddress(ctx, request, email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) updateCustomerCarHandler(c echo.Context) error {
	ctx := c.Request().Context()
	car_id := c.QueryParam("_id")

	request := new(models.UpdateCustomerCar)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	if err := cr.CustomerServices.UpdateCustomerCar(ctx, request, car_id); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) deleteCustomerAddressHandler(c echo.Context) error {
	ctx := c.Request().Context()
	address_id := c.QueryParam("_id")

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	if err := cr.CustomerServices.DeleteCustomerAddress(ctx, address_id, email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) updateCustomerDefaultAddressHandler(c echo.Context) error {
	ctx := c.Request().Context()
	address_id := c.QueryParam("_id")

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	if err := cr.CustomerServices.UpdateCustomerDefaultAddress(ctx, address_id, email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) updateCustomerAddressHandler(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	address_id := c.QueryParam("_id")
	request := new(models.RegisterAddressRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	location_name := cr.CustomerServices.CheckExistLocationName(ctx, email, request.LocationName)
	if location_name != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This location name is already exists"})
	}

	if err := cr.CustomerServices.UpdateCustomerAddress(ctx, request, email, address_id); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) deleteCustomerCarHandler(c echo.Context) error {
	ctx := c.Request().Context()
	car_id := c.QueryParam("_id")

	if err := cr.CustomerServices.DeleteCustomerCar(ctx, car_id); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) updateCustomerDefaultCarHandler(c echo.Context) error {
	ctx := c.Request().Context()
	address_id := c.QueryParam("_id")

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	if err := cr.CustomerServices.UpdateCustomerDefaultCar(ctx, address_id, email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) customerFavoriteAddressHandler(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "error to get user profile email"})
	}

	area := cr.CustomerServices.CustomerFavoriteArea(ctx, account.FavoriteArea)
	if area == nil {
		return c.JSON(http.StatusOK, models.SearchAreaResponse{Data: nil, Status: http.StatusOK})
	}
	return c.JSON(http.StatusOK, models.SearchAreaResponse{Data: area, Status: http.StatusOK})
}

func (cr CustomerRouters) updateCustomerFavoriteCarHandler(c echo.Context) error {
	ctx := c.Request().Context()
	parking_id := c.QueryParam("parking_id")
	action := c.QueryParam("action")

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	if err := cr.CustomerServices.UpdateCustomerFavoriteArea(ctx, email, parking_id, action); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) searchLocationParkingAreaHandler(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(models.SearchQueryRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	if request.MinStart != 0 && request.MinStart != 30 || request.MinEnd != 0 && request.MinEnd != 30 {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Min interval should be 0 or 30"})
	}
	areas := cr.SearchServices.Search(ctx, request, request.Keyword)
	if areas == nil {
		return c.JSON(http.StatusNoContent, models.MessageResponse{Message: "No matching areas found"})
	}

	return c.JSON(http.StatusOK, models.SearchAreaResponse{Data: areas, Status: http.StatusOK})
}

func (cr CustomerRouters) customerRewardHandler(c echo.Context) error {
	ctx := c.Request().Context()

	reward := cr.CustomerServices.CustomerReward(ctx)

	return c.JSON(http.StatusOK, models.RewardResponse{Data: reward, Status: http.StatusOK})

}

func (cr CustomerRouters) customerRewardDetail(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("_id")

	reward := cr.CustomerServices.CustomerRewardDetail(ctx, id)

	return c.JSON(http.StatusOK, models.RewardDetailResponse{Data: reward, Status: http.StatusOK})

}

func (cr CustomerRouters) customerRedeemReward(c echo.Context) error {
	ctx := c.Request().Context()
	reward_id := c.QueryParam("_id")
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Not found Account ID"})
	}
	code, err := cr.CustomerServices.CustomerRedeemReward(ctx, reward_id, *account)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: code})

}

func (cr CustomerRouters) customerMyRedeemReward(c echo.Context) error {
	ctx := c.Request().Context()
	reward_id := c.QueryParam("_id")
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Not found Account ID"})
	}

	if reward_id == "" {
		var customer_rewards []models.CustomerRedeemReward
		for _, reward := range account.Reward {
			currentTime := time.Now()
			if reward.ExpiredDate.After(currentTime) {
				duration := reward.ExpiredDate.Sub(currentTime)
				hours := duration.Hours()
				reward_detail := cr.CustomerServices.CustomerRewardDetail(ctx, reward.ID.Hex())
				data := models.CustomerRedeemReward{
					ID:              reward.ID,
					Name:            reward_detail.Name,
					Title:           reward_detail.Title,
					Description:     reward_detail.Description,
					PreviewImageURL: reward_detail.PreviewImageURL,
					Webhook:         reward_detail.Webhook,
					Condition:       reward_detail.Condition,
					QuotaCount:      reward_detail.QuotaCount,
					CreateBy:        reward_detail.CreateBy,
					BarcodeURL:      reward.BarcodeURL,
					ExpiredHour:     int(hours),
				}

				customer_rewards = append(customer_rewards, data)
			}
		}
		return c.JSON(http.StatusOK, models.CutomserRewardList{Data: customer_rewards, Status: 200})
	} else {
		var customer_rewards []models.CustomerRedeemReward
		for _, reward := range account.Reward {
			if reward.ID.Hex() == reward_id {
				currentTime := time.Now()
				if reward.ExpiredDate.After(currentTime) {
					duration := reward.ExpiredDate.Sub(currentTime)
					hours := duration.Hours()
					reward_detail := cr.CustomerServices.CustomerRewardDetail(ctx, reward.ID.Hex())
					data := models.CustomerRedeemReward{
						ID:              reward.ID,
						Name:            reward_detail.Name,
						Title:           reward_detail.Title,
						Description:     reward_detail.Description,
						PreviewImageURL: reward_detail.PreviewImageURL,
						Webhook:         reward_detail.Webhook,
						Condition:       reward_detail.Condition,
						QuotaCount:      reward_detail.QuotaCount,
						CreateBy:        reward_detail.CreateBy,
						BarcodeURL:      reward.BarcodeURL,
						ExpiredHour:     int(hours),
					}

					customer_rewards = append(customer_rewards, data)
				}
				return c.JSON(http.StatusOK, models.CutomserRewardList{Data: customer_rewards, Status: 200})
			}
		}
	}

	return c.JSON(http.StatusOK, models.CutomserRewardList{Data: nil, Status: 200})

}

func (cr CustomerRouters) MyHistoryPoint(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)
	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Not found Account ID"})
	}

	var paids []models.CutomserHistoryPoint
	loc, _ := time.LoadLocation("Asia/Bangkok")
	for _, reward := range account.Reward {
		t := reward.TimeStamp.In(loc)
		TimeStampString := utility.FormatThaiDateTime(t)
		paid := models.CutomserHistoryPoint{
			Content:         fmt.Sprintf("คุณได้แลกคูปอง %s", reward.Name),
			Type:            "paid",
			Point:           reward.Point,
			TimeStamp:       reward.TimeStamp,
			TimeStampString: TimeStampString,
		}
		paids = append(paids, paid)

	}

	received, err := cr.ReserveService.MyReservePaymentComplete(ctx, email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Not found reserve history"})
	}

	received = append(received, paids...)
	sort.Slice(received, func(i, j int) bool {
		return received[i].TimeStamp.Before(received[j].TimeStamp)
	})
	return c.JSON(http.StatusOK, models.CutomserHistoryPointResponse{Data: received, Status: 200})

}
