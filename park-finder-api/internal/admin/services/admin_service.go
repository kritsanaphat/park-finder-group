package services

import (
	"context"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as AdminServices) AdminGetParkingArea(ctx context.Context, status string) *[]models.ParkingArea {

	area := as.ParkingAreaStorage.FindParkingAreaByStatus(ctx, status)
	if area == nil {
		return nil
	}

	return area
}

func (as AdminServices) AdminGetTransaction(ctx context.Context, year, month string) (*[]models.AdminTransactionResponse, error) {
	var all_data []models.AdminTransactionResponse
	all_account := as.ProviderStorage.FindAllListProviderID(ctx)
	for _, account := range all_account {
		list_parking_id := as.ParkingAreaStorage.FindListParkingIDByProviderID(ctx, account.ID)
		if list_parking_id != nil {
			sum, err := as.TransactionStorage.FindTransactionFromAdmin(ctx, list_parking_id, month, year)
			if err != nil {
				return nil, nil
			}
			is_pay, price := as.ReceiptStorage.FindExistReceipt(ctx, month, year, account.ID, sum)
			if price != 0 {
				sum = price
			}
			data := models.AdminTransactionResponse{
				BankAccount: account.BankAccount,
				Count:       len(list_parking_id),
				Sum:         sum,
				ProviderID:  account.ID,
				IsPay:       is_pay,
			}
			all_data = append(all_data, data)
		}

	}

	return &all_data, nil
}

func (as AdminServices) AdminSubmitReceipt(ctx context.Context, image_url, provider_id, month, year string, price int) error {
	_id, err := primitive.ObjectIDFromHex(provider_id)
	if err != nil {
		return err
	}
	data := models.Receipt{
		ID:              primitive.NewObjectID(),
		ProviderID:      _id,
		Year:            year,
		Month:           month,
		ReceiptImageUrl: image_url,
		Price:           price,
		TimeStamp:       time.Now(),
	}

	_, err = as.ReceiptStorage.InsertLog(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (as AdminServices) AdminUpdateParkingArea(ctx context.Context, parking_id, status, description string) error {

	_, err := as.ParkingAreaStorage.UpdateParkingAreaByStatus(ctx, parking_id, status, description)
	if err != nil {
		return err
	}

	return nil
}
func (as AdminServices) AdminCheckParkingAreaDetail(ctx context.Context, parking_id string) *models.ParkingArea {

	parkings := new(models.ParkingArea)
	err := as.ParkingAreaStorage.FindParkingByParkingID(ctx, parking_id, parkings)
	if err != nil {
		return nil
	}

	return parkings
}

func (as AdminServices) AddReward(ctx context.Context, user *models.AddRewardRequest, email string) (string, error) {

	layout := "02/01/2006,15.04"
	expiredDate, err := time.Parse(layout, user.ExpiredDate)
	if err != nil {
		return "", err
	}
	_id := primitive.NewObjectID()
	data := &models.Reward{
		ID:              _id,
		Name:            user.Name,
		Point:           user.Point,
		Title:           user.Title,
		Description:     user.Description,
		PreviewImageURL: user.PreviewImageURL,
		Webhook:         user.Webhook,
		ExpiredDate:     expiredDate,
		TimeStamp:       time.Now(),
		Condition:       user.Condition,
		QuotaCount:      user.QuotaCount,
		CreateBy:        email,
	}
	_, err = as.RewardStorage.InsertReward(ctx, data)
	if err != nil {
		return "", err
	}

	return _id.Hex(), nil
}
