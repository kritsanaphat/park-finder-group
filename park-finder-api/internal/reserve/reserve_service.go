package reserve

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (rs ReserveServices) ReserveParking(ctx context.Context, email, order_id string, req *models.ReserveRequest) (string, error, *primitive.ObjectID) {

	id := primitive.NewObjectID()
	parking_id, err := primitive.ObjectIDFromHex(req.ParkingID)
	if err != nil {
		fmt.Println("Error parking_id not incorrect format:", err)
	}
	car_id, err := primitive.ObjectIDFromHex(req.CarID)
	if err != nil {
		fmt.Println("Error parking_id not incorrect format:", err)
	}
	provider_id, err := primitive.ObjectIDFromHex(req.ProviderID)
	if err != nil {
		fmt.Println("Error parking_id not incorrect format:", err)
	}

	area := new(models.ParkingArea)
	err = rs.ParkingAreaStorage.FindParkingByParkingID(ctx, req.ParkingID, area)
	if err != nil {
		fmt.Println("Error parking_id not correct format:", err)
	}

	data := &models.Reservation{
		ID:                id,
		CustomerEmail:     email,
		ProviderID:        provider_id,
		ParkingID:         parking_id,
		CarID:             car_id,
		Status:            "Pending",
		OrderID:           order_id,
		DateStart:         req.DateStart,
		DateEnd:           req.DateEnd,
		Type:              req.Type,
		HourStart:         req.HourStart,
		MinStart:          req.MinStart,
		HourEnd:           req.HourEnd,
		MinEnd:            req.MinEnd,
		PaymentChanel:     req.PaymentChanel,
		TimeStamp:         time.Now(),
		Price:             req.Price,
		Address_Full:      area.Address.AddressText + " " + area.Address.Sub_district + " " + area.Address.District + " " + area.Address.Province + " " + area.Address.Postal_code,
		ParkingPictureUrl: area.ParkingPictureUrl,
		ParkingName:       area.ParkingName,
		ModuleCode:        "",
		IsExtend:          false,
	}

	_, err = rs.ReserveStorage.InsertLogReservation(ctx, *data)
	if err != nil {
		return "", err, nil
	}

	go rs.HttpClient.SendCancelReserveAPI(order_id)

	return order_id, nil, &id
}

func (rs ReserveServices) ReserveParkingInAdvance(ctx context.Context, email, order_id string, req *models.ReserveRequest) (string, error, *primitive.ObjectID) {

	id := primitive.NewObjectID()
	parking_id, err := primitive.ObjectIDFromHex(req.ParkingID)
	if err != nil {
		fmt.Println("Error parking_id not incorrect format:", err)
	}
	car_id, err := primitive.ObjectIDFromHex(req.CarID)
	if err != nil {
		fmt.Println("Error parking_id not incorrect format:", err)
	}
	provider_id, err := primitive.ObjectIDFromHex(req.ProviderID)
	if err != nil {
		fmt.Println("Error parking_id not incorrect format:", err)
	}

	area := new(models.ParkingArea)
	err = rs.ParkingAreaStorage.FindParkingByParkingID(ctx, req.ParkingID, area)
	if err != nil {
		fmt.Println("Error parking_id not correct format:", err)
	}

	status := "Pending Approval"

	data := &models.Reservation{
		ID:                id,
		CustomerEmail:     email,
		ProviderID:        provider_id,
		ParkingID:         parking_id,
		CarID:             car_id,
		Status:            status,
		OrderID:           order_id,
		DateStart:         req.DateStart,
		DateEnd:           req.DateEnd,
		Type:              req.Type,
		HourStart:         req.HourStart,
		MinStart:          req.MinStart,
		HourEnd:           req.HourEnd,
		MinEnd:            req.MinEnd,
		PaymentChanel:     req.PaymentChanel,
		TimeStamp:         time.Now(),
		Price:             req.Price,
		Address_Full:      area.Address.AddressText + " " + area.Address.Sub_district + " " + area.Address.District + " " + area.Address.Province + " " + area.Address.Postal_code,
		ParkingPictureUrl: area.ParkingPictureUrl,
		ParkingName:       area.ParkingName,
		ModuleCode:        "",
	}

	_, err = rs.ReserveStorage.InsertLog(ctx, data)
	if err != nil {
		return "", err, nil
	}
	go rs.HttpClient.SendCancelReserveAPI(order_id)

	return order_id, nil, &id
}

func (rs ReserveServices) MyReserve(ctx context.Context, email, parking_id, status string) ([]models.Reservation, error) {

	if status == "on_working" {
		reserve, err := rs.ReserveStorage.FindProcessReservationByCustomerEmail(ctx, email, parking_id)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return reserve, nil

	} else if status == "fail" {
		reserve, err := rs.ReserveStorage.FindCancelReservationByCustomerEmail(ctx, email, parking_id)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return reserve, nil

	} else if status == "successful" {
		reserve, err := rs.ReserveStorage.FindSuccessfulReservationByCustomerEmail(ctx, email, parking_id)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}

		return reserve, nil

	}

	return nil, nil
}

func (rs ReserveServices) MyReserveProvider(ctx context.Context, provider_id primitive.ObjectID, parking_id, status string) ([]models.Reservation, error) {

	if status == "on_working" {
		reserve, err := rs.ReserveStorage.FindProcessReservationByProviderID(ctx, provider_id, parking_id)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return reserve, nil

	} else if status == "fail" {
		reserve, err := rs.ReserveStorage.FindCancelReservationByProviderID(ctx, provider_id, parking_id)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return reserve, nil

	} else if status == "successful" {
		reserve, err := rs.ReserveStorage.FindSuccessfulReservationByProviderID(ctx, provider_id, parking_id)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}

		return reserve, nil

	}

	return nil, nil
}

func (rs ReserveServices) StartReserveParking(ctx context.Context, customer_id primitive.ObjectID, order_id, module_code, parking_id, license_plate, action, parking_name string) (error, bool) {
	is_found := rs.HttpClient.SendCameraServiceToDectectIncommingCar(module_code, license_plate)
	if action == "force" || is_found {
		fmt.Println(order_id)
		log := rs.FindReserveDetailByOrderID(ctx, order_id)
		if log == nil {
			return errors.New("Order id doesn't exist"), false
		}
		parking_Id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error parking_id not incorrect format:", err)
			return err, false
		}

		_, err = rs.ReserveStorage.UpdateLogByInterface(ctx, bson.M{
			"order_id":   order_id,
			"status":     "Process",
			"parking_id": parking_Id,
		}, bson.M{"$set": bson.M{"module_code": module_code, "status": "Parking"}})
		if err != nil {
			return err, false
		}

	} else {
		reserve := rs.ReserveStorage.FindReservationByOrderID(ctx, order_id)
		if reserve == nil {
			return errors.New("Order id doesn't exist"), false
		}
		pic_url, err := rs.HttpClient.SendCameraServiceTCaptureCar(module_code)
		if err != nil {
			return errors.New("Error form capture image"), false
		}
		go rs.NotificationService.VertifyCustomerCarNotification(ctx, customer_id, license_plate, module_code, pic_url, reserve)

	}

	return nil, is_found
}

func (rs ReserveServices) CheckExistOrderID(ctx context.Context, order_id string) *models.Reservation {
	reserve := rs.ReserveStorage.FindReservationByOrderID(ctx, order_id)

	return reserve
}

func (rs ReserveServices) CheckOrderIDByParkingIDAndCustomerEmail(ctx context.Context, provider_id, customer_email string) *models.Reservation {
	reserve := rs.ReserveStorage.FindReservationIDByParkingIDAndCustomerEmail(ctx, provider_id, customer_email)
	return reserve
}

func (rs ReserveServices) FindParkingDetail(ctx context.Context, parking_id string) *models.ParkingArea {

	area := new(models.ParkingArea)
	err := rs.ParkingAreaStorage.FindParkingByParkingID(ctx, parking_id, area)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return area
}

func (rs ReserveServices) FindReserveDetail(ctx context.Context, reservation_id string) *models.Reservation {

	reservation := rs.ReserveStorage.FindReservationByReserveID(ctx, reservation_id)
	if reservation == nil {
		return nil
	}
	return reservation
}

func (rs ReserveServices) FindReserveDetailByOrderID(ctx context.Context, order_id string) *models.Reservation {

	reservation := rs.ReserveStorage.FindReservationByOrderID(ctx, order_id)
	if reservation == nil {
		return nil
	}
	return reservation
}

func (rs ReserveServices) CreateReview(ctx context.Context, customer_id, fn, ln, comment, parking_id, order_id string, review_score int) error {
	result, err := rs.ParkingAreaStorage.PushReview(ctx, customer_id, fn, ln, comment, parking_id, order_id, review_score)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}
	return nil
}

func (rs ReserveServices) CheckReserveInNextHour(ctx context.Context, parking_id primitive.ObjectID, hour_start, hour_end int, date_end string) (*[]models.Reservation, error) {
	var reservation []models.Reservation
	parking_detail := new(models.ParkingArea)
	err := rs.ParkingAreaStorage.FindParkingByParkingID(ctx, parking_id.Hex(), parking_detail)
	if err != nil {
		return nil, err
	}
	reserves, err := rs.ReserveStorage.FindIncommingReservationByParkingID(ctx, parking_id)
	if err != nil {
		return nil, err
	}
	fmt.Println("Incomming Reserve ", reserves)
	for _, reserve := range reserves {
		if _isOverlapDate(date_end, reserve.DateStart, reserve.DateEnd) {
			if hour_end >= reserve.HourStart {
				reservation = append(reservation, reserve)

			}
		}
	}
	slot_parking := parking_detail.ToatalParkingCount - 1
	for slot_parking > 0 {

		reservation = _popFirstElement(reservation)
		slot_parking -= 1
		if len(reservation) == 0 {
			break
		}

	}

	return &reservation, nil
}

func (rs ReserveServices) ConfirmReservationInAdvance(ctx context.Context, order_id string) error {
	if err := rs.ReserveStorage.UpdateConfirmStatusToProcessByOrderID(ctx, order_id); err != nil {
		return err
	}
	return nil
}

func (rs ReserveServices) ExtendReserve(ctx context.Context, order_id, action, date_end string, hour_end, min_end int) error {
	if action == "normal" {
		_, err := rs.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"hour_end": hour_end, "status": "Pending", "is_extend": true}})
		if err != nil {
			return err
		}
		go rs.HttpClient.SendRemoveJobAPI("BTOR_" + order_id)
		go rs.HttpClient.SendRemoveJobAPI("TOR_" + order_id)
		go rs.HttpClient.SendRemoveJobAPI("ATOR_" + order_id)
		go rs.HttpClient.SendCancelExtendReserveAPI(order_id, hour_end)
	} else if action == "automatic" {
		_, err := rs.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"hour_end": hour_end, "status": "Parking", "is_extend": true}})
		if err != nil {
			return err
		}
		go rs.HttpClient.SendRemoveJobAPI("BTOR_" + order_id)
		go rs.HttpClient.SendRemoveJobAPI("TOR_" + order_id)
		go rs.HttpClient.SendRemoveJobAPI("ATOR_" + order_id)
		go rs.HttpClient.SendTimeoutReserveAPI(order_id, date_end, strconv.Itoa(hour_end), strconv.Itoa(min_end))
	}
	return nil
}

func (rs ReserveServices) CaptureCarReserve(ctx context.Context, receiver_id primitive.ObjectID, module_code, parking_name string) (string, error) {
	pic_url, err := rs.HttpClient.SendCameraServiceTCaptureCar(module_code)
	if err != nil {
		return "", err
	}

	return pic_url, nil
}

func (rs ReserveServices) UpdateReserveStatus(ctx context.Context, status, order_id string) error {
	_, err := rs.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}

	return nil
}

func (rs ReserveServices) UpdateReserveStatusAndRemoveJob(ctx context.Context, order_id string) error {
	_, err := rs.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"status": "Cancel"}})
	if err != nil {
		return err
	}
	go rs.HttpClient.SendRemoveJobAPI("BTOR_" + order_id)
	go rs.HttpClient.SendRemoveJobAPI("TOR_" + order_id)
	go rs.HttpClient.SendRemoveJobAPI("ATOR_" + order_id)

	return nil
}

func (rs ReserveServices) MyReservePaymentComplete(ctx context.Context, email string) ([]models.CutomserHistoryPoint, error) {

	reserves, err := rs.ReserveStorage.FindComplePaymentReservationByCustomerEmail(ctx, email)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var data []models.CutomserHistoryPoint
	for _, reserve := range reserves {

		StartDate := utility.FormatThaiDateTimeFromString(reserve.DateStart)

		EndDate := utility.FormatThaiDateTimeFromString(reserve.DateEnd)

		temp := models.CutomserHistoryPoint{
			Content: fmt.Sprintf("คุณได้จองที่จอดรถ %s", reserve.ParkingName),
			Point:   reserve.Price,
			Type:    "received",
			TimeStampString: fmt.Sprintf("%s เวลา %d:%d ถึง %s เวลา %d:%d ",
				StartDate,
				reserve.HourStart,
				reserve.MinStart,
				EndDate,
				reserve.HourEnd,
				reserve.MinEnd,
			),
			TimeStamp: reserve.TimeStamp,
		}
		data = append(data, temp)
	}
	return data, nil

}

func (rs ReserveServices) ReportReservation(ctx context.Context, customer_id, provider_id primitive.ObjectID, content, order_id string) error {
	data := &models.Report{
		ID:         primitive.NewObjectID(),
		CustomerID: customer_id,
		ProviderID: provider_id,
		OrderID:    order_id,
		Content:    content,
		TimeStamp:  time.Now(),
	}
	_, err := rs.ReportStorage.InsertLog(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// if overlap return true
func _isOverlap(range1Start, range1End, range2Start, range2End int) bool {
	return range1Start >= range2Start && range1End <= range2End
}

func _isOverlapDate(date_end, reserve_date_start, reserve_date_end string) bool {
	input, _ := time.Parse("2006-01-02", date_end)

	r_date_start, _ := time.Parse("2006-01-02", reserve_date_start)

	r_date_end, _ := time.Parse("2006-01-02", reserve_date_end)

	if input.Equal(r_date_start) || input.After(r_date_start) && (input.Before(r_date_end) || input.Equal(r_date_end)) {
		return true
	} else {
		return false
	}
}

func _popFirstElement(slice []models.Reservation) []models.Reservation {
	if len(slice) == 0 {
		return slice
	}
	copy(slice, slice[1:])

	return slice[:len(slice)-1]
}

func _diffDate(end string) int {
	layout := "2006-01-02"
	currentTime := time.Now()

	date1, _ := time.Parse(layout, currentTime.Format(layout))
	date2, _ := time.Parse(layout, end)

	days := date2.Sub(date1).Hours() / 24
	return int(days)
}
