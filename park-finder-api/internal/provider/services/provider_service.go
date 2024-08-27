package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ps ProviderServices) UpdateProviderProfile(ctx context.Context, email string, user *models.UpdateProfileRequest) error {
	location, _ := time.LoadLocation("Asia/Bangkok")
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": models.UpdateCustomerProfile{
			SSN:               user.SSN,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Birthday:          user.Birthday,
			Phone:             user.Phone,
			ProfilePictureURL: user.ProfilePictureURL,
			TimeStamp:         time.Now().In(location),
		}}
	result, err := ps.ProviderAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (ps ProviderServices) UpdateProviderBankAccount(ctx context.Context, email string, bank *models.BankAccount) error {
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": bson.M{
			"bank_account": bson.M{
				"account_book_image_url": bank.AccountBookImageUrl,
				"bank_name":              bank.BankName,
				"account_name":           bank.AccountName,
				"account_number":         bank.AccountNumber,
			},
		},
	}
	result, err := ps.ProviderAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (ps ProviderServices) ProviderRegisterAreaLocaion(ctx context.Context, area *models.RegisterParkingAreaFirstStepRequest, email string) (interface{}, error) {

	filter := bson.M{"email": email}

	user := new(models.ProviderAccount)
	ps.ProviderAccoutStorage.FindAccountInterface(ctx, filter, user)

	open_detail := models.Daily{}

	address := models.ParkingAddress{
		AddressText:  area.AddressText,
		Sub_district: area.Sub_district,
		District:     area.District,
		Province:     area.Province,
		Postal_code:  area.Postal_code,
		Latitude:     area.Latitude,
		Longitude:    area.Longitude,
	}

	review := []models.ReviewParkingArea{}
	data := &models.ParkingArea{
		ID:                     primitive.NewObjectID(),
		ProviderID:             user.ID,
		ParkingName:            area.ParkingName,
		OpenDetail:             open_detail,
		Address:                address,
		Tag:                    area.Tag,
		OpenStatus:             false,
		Review:                 review,
		TimeStamp:              time.Now(),
		StatusApplyDescription: "",
		StatusApply:            "wait for document",
		ReserveLog:             []models.ReserveLog{},
		Distance:               0.0,
		TimeStampClose:         nil,
	}

	id, err := ps.ParkingAreaStorage.InsertParkingArea(ctx, data)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (ps ProviderServices) ProviderRegisterAreaDocument(ctx context.Context, area *models.RegisterParkingAreaDocumentStepRequest, email, parking_id string) error {
	user := new(models.ProviderAccount)
	filter := bson.M{"email": email}

	ps.ProviderAccoutStorage.FindAccountInterface(ctx, filter, user)

	result, err := ps.ParkingAreaStorage.UpdateOpenAreaRegisterDocumentClose(ctx, user.ID, *area, parking_id)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no matching document found")
	}

	return nil
}

func (ps ProviderServices) GetProviderArea(ctx context.Context, id string) []models.ParkingArea {

	var areas []models.ParkingArea
	areas, err := ps.ParkingAreaStorage.FindParkingByProviderID(ctx, id)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return areas
}

func (ps ProviderServices) UpdateOpenAreaDailyStatus(ctx context.Context, daily *models.UpdateOpenAreaDailyStatusRequest) error {

	result, err := ps.ParkingAreaStorage.UpdateOpenAreaDailyStatus(ctx, *daily)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (ps ProviderServices) UpdateOpenAreaQuickStatus(ctx context.Context, id, status string, range_time int) error {

	if status == "close" {
		result, err := ps.ParkingAreaStorage.UpdateOpenAreaQuickStatusClose(ctx, id, range_time)
		if err != nil {
			return err
		}
		if result.ModifiedCount == 0 {
			return errors.New("no documents were modified")
		}
		go ps.HttpClient.SendOpenStatusArea(id, range_time)
	}

	if status == "open" {
		result, err := ps.ParkingAreaStorage.UpdateOpenAreaQuickStatusOpen(ctx, id, range_time)
		if err != nil {
			return err
		}
		if result.ModifiedCount == 0 {
			return errors.New("no documents were modified")
		}
		go ps.HttpClient.SendRemoveJobAPI("OPA" + id)

	}

	return nil
}

func (ps ProviderServices) UpdateOpenAreaInAdvanceStatus(ctx context.Context, id, status string, date []string) error {
	if status == "close" {
		result, err := ps.ParkingAreaStorage.PushOpenAreaDateClose(ctx, id, date)
		if err != nil {
			return err
		}
		if result.ModifiedCount == 0 {
			return errors.New("no documents were modified")
		}
	}

	if status == "open" {
		result, err := ps.ParkingAreaStorage.PullOpenAreaDateClose(ctx, id, date)
		if err != nil {
			return err
		}
		if result.ModifiedCount == 0 {
			return errors.New("no documents were modified")
		}
	}
	return nil
}

func (ps ProviderServices) GetProviderProfitDaily(ctx context.Context, list_parking_id []string) (*models.DailyProfitResponse, int, error) {
	date, sum, count, err := ps.TransactionStorage.FindProfitDaily(ctx, list_parking_id)
	if err != nil {
		return nil, 0, err
	}

	response := models.DailyProfitResponse{
		Date: date,
		Sum:  sum,
	}

	return &response, count, nil
}

func (ps ProviderServices) GetProviderProfitWeekly(ctx context.Context, list_parking_id []string) (*[]models.DailyProfitResponse, int, error) {
	sum, count, err := ps.TransactionStorage.FindProfitWeekly(ctx, list_parking_id)
	if err != nil {
		return nil, 0, err
	}

	return &sum, count, nil
}

func (ps ProviderServices) GetProviderProfitMontly(ctx context.Context, list_parking_id []string) (*[]models.DailyProfitResponse, int, error) {
	sum, count, err := ps.TransactionStorage.FindProfitMontly(ctx, list_parking_id)
	if err != nil {
		return nil, 0, err
	}

	return &sum, count, nil
}

func (ps ProviderServices) GetProviderProfitYearly(ctx context.Context, list_parking_id []string) (*[]models.DailyProfitResponse, int, error) {
	sum, count, err := ps.TransactionStorage.FindProfitYearly(ctx, list_parking_id)
	if err != nil {
		return nil, 0, err
	}

	return &sum, count, nil
}
func (ps ProviderServices) CheckValidProviderArea(ctx context.Context, provider_id, parking_area_id string) error {

	err := ps.ParkingAreaStorage.FindParkingAreaByProviderID(ctx, provider_id, parking_area_id)
	if err == mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func (ps ProviderServices) CheckQuickAvailabilityCloseProviderArea(ctx context.Context, parking_area_id string, range_time int) *[]models.Reservation {

	reservation, err := ps.ReserveStorage.FindActiveReservationByParkingID(ctx, parking_area_id)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading location:", err)
	}

	currentTime := time.Now().In(location)
	var reservations []models.Reservation

	hour_start := currentTime.Hour()
	hour_end := currentTime.Add(time.Duration(range_time) * time.Hour).Hour()

	for _, reserve := range reservation {
		if !_isOverlap(hour_start, hour_end, reserve.HourStart, reserve.HourEnd) {
			reservations = append(reservations, reserve)
		}
	}
	if len((reservation)) == 0 {
		return nil
	}
	return &reservations
}

func (ps ProviderServices) CheckDailyAvailabilityCloseProviderArea(ctx context.Context, request models.UpdateOpenAreaDailyStatusRequest) *[]models.Reservation {

	reservation, err := ps.ReserveStorage.FindActiveReservationByParkingID(ctx, request.ParkingAreaID)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	var reservations []models.Reservation

	for _, reserve := range reservation {
		daily_list := utility.ConvertToDayNames(reserve.DateStart, reserve.DateEnd)
		for _, daily := range daily_list {
			if !_checkCanColseByDailyAndHour(daily, reserve.HourStart, reserve.HourEnd, request) {
				if reserve.Status == "Process" {
					reservations = append(reservations, reserve)
					break
				}
			}
		}
	}

	if len((reservation)) == 0 {
		return nil
	}
	return &reservations
}

func (ps ProviderServices) CheckUpdatePriceArea(ctx context.Context, parking_area_id string, price int16) error {
	_id, err := primitive.ObjectIDFromHex(parking_area_id)
	if err != nil {
		return err
	}
	reservation, err := ps.ReserveStorage.FindStillParkingReservationByParkingID(ctx, _id)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if reservation != nil {
		return errors.New("Can't Change Price Because have a customer sitll parking")
	}
	_, err = ps.ParkingAreaStorage.UpdateParkingAreaByPrice(ctx, _id, price)
	if err != nil {
		return errors.New("Can't Change Price")
	}
	return nil
}

func (ps ProviderServices) FineProvider(ctx context.Context, _id primitive.ObjectID, fine int) error {
	err := ps.ProviderAccoutStorage.UpdateFineProvider(ctx, _id, fine)
	if err != nil {
		return err
	}
	return nil
}

func _checkCanColseByDailyAndHour(daily string, hour_start, hour_end int, request models.UpdateOpenAreaDailyStatusRequest) bool {
	switch daily {
	case "Monday":
		if _isOverlap(hour_start, hour_end, request.Monday.OpenTime, request.Monday.CloseTime) {
			return true
		}
	case "Tuesday":
		if _isOverlap(hour_start, hour_end, request.Tuesday.OpenTime, request.Tuesday.CloseTime) {
			return true
		}
	case "Wednesday":
		if _isOverlap(hour_start, hour_end, request.Wednesday.OpenTime, request.Wednesday.CloseTime) {
			return true
		}
	case "Thursday":
		if _isOverlap(hour_start, hour_end, request.Thursday.OpenTime, request.Thursday.CloseTime) {
			return true
		}
	case "Friday":
		if _isOverlap(hour_start, hour_end, request.Friday.OpenTime, request.Friday.CloseTime) {
			return true
		}
	case "Saturday":
		if _isOverlap(hour_start, hour_end, request.Saturday.OpenTime, request.Saturday.CloseTime) {
			return true
		}
	case "Sunday":
		if _isOverlap(hour_start, hour_end, request.Sunday.OpenTime, request.Sunday.CloseTime) {
			return true
		}
	}
	return false
}

func _isOverlap(range1Start, range1End, range2Start, range2End int) bool {
	fmt.Println(range1Start, range1End)
	return range1Start >= range2Start && range1End <= range2End
}
