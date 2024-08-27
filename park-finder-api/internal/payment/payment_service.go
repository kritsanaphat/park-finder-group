package payment

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ps PaymentServices) LinePayReserve(ctx context.Context, email, order_id, parking_id, action string, req *models.LineReserveRequest) (string, error) {
	if action == "fine" {

		url_call_back := "https://park-finder.online/webhook/line-pay/reserve/callback/fine"
		parking_Id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error parking_id not incorrect format:", err)
		}
		url := "/v3/payments/request"
		var pk []models.Package
		if req.CashBack == 0 {
			pk = []models.Package{
				{
					PackageID: req.ProviderID,
					Amount:    req.Quantity * float32(req.Price),
					Name:      req.ParkingName,
					Products: []models.Product{
						{
							Name:     req.ParkingName,
							Quantity: req.Quantity,
							Price:    req.Price,
							ImageURL: "https://pay-store.line.com/images/pen_brown.jpg",
						},
					},
				},
			}
		} else {
			pk = []models.Package{
				{
					PackageID: req.ProviderID,
					Amount:    req.Quantity*float32(req.Price) - float32(req.CashBack),
					Name:      req.ParkingName,
					Products: []models.Product{
						{
							Name:     req.ParkingName,
							Quantity: 1,
							Price:    int(req.Quantity*float32(req.Price) - float32(req.CashBack)),
							ImageURL: "https://pay-store.line.com/images/pen_brown.jpg",
						},
					},
				},
			}
		}

		request := models.LineReserveAPIRequest{
			Amount:   (req.Quantity * float32(req.Price)) - float32(req.CashBack),
			Currency: "THB",
			OrderId:  "FINE_" + order_id,
			Packages: pk,
			RedirectUrls: models.RedirectUrls{
				ConfirmUrl:     url_call_back,
				CancelUrl:      os.Getenv("CANCEL_LINE_RESERVE"),
				ConfirmUrlType: "SERVER",
			},
		}

		res := new(models.LineReserveAPIResponse)
		err = ps.HttpClient.SendLinePayRequestAPI(request, url, res)
		if err != nil {
			return "", err
		}

		InfoReserve := models.InfoReserve{
			PaymentUrl: models.PaymentUrl{
				Web: res.Info.PaymentUrl.Web,
				App: res.Info.PaymentUrl.App,
			},
			PaymentAccessToken: res.Info.PaymentAccessToken,
			TransactionId:      res.Info.TransactionId,
		}

		data := &models.TransactionLinePay{
			ID:            primitive.NewObjectID(),
			ParkingID:     parking_Id,
			CustomerEmail: email,
			OrderId:       "FINE_" + order_id,
			Packages:      pk,
			Info:          InfoReserve,
			Status:        "Pending",
			TimeStamp:     time.Now(),
		}

		_, err = ps.TransactionStorage.InsertLog(ctx, data)
		if err != nil {
			return "", err
		}
		return res.Info.PaymentUrl.Web, nil

	} else {
		url_call_back := "https://park-finder.online/webhook/line-pay/reserve/callback"
		parking_Id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error parking_id not incorrect format:", err)
		}
		url := "/v3/payments/request"
		var pk []models.Package
		if req.CashBack == 0 {
			pk = []models.Package{
				{
					PackageID: req.ProviderID,
					Amount:    req.Quantity * float32(req.Price),
					Name:      req.ParkingName,
					Products: []models.Product{
						{
							Name:     req.ParkingName,
							Quantity: req.Quantity,
							Price:    req.Price,
							ImageURL: "https://pay-store.line.com/images/pen_brown.jpg",
						},
					},
				},
			}
		} else {
			pk = []models.Package{
				{
					PackageID: req.ProviderID,
					Amount:    req.Quantity*float32(req.Price) - float32(req.CashBack),
					Name:      req.ParkingName,
					Products: []models.Product{
						{
							Name:     req.ParkingName,
							Quantity: 1,
							Price:    int(req.Quantity*float32(req.Price) - float32(req.CashBack)),
							ImageURL: "https://pay-store.line.com/images/pen_brown.jpg",
						},
					},
				},
			}
		}

		request := models.LineReserveAPIRequest{
			Amount:   (req.Quantity * float32(req.Price)) - float32(req.CashBack),
			Currency: "THB",
			OrderId:  order_id,
			Packages: pk,
			RedirectUrls: models.RedirectUrls{
				ConfirmUrl:     url_call_back,
				CancelUrl:      os.Getenv("CANCEL_LINE_RESERVE"),
				ConfirmUrlType: "SERVER",
			},
		}

		res := new(models.LineReserveAPIResponse)
		err = ps.HttpClient.SendLinePayRequestAPI(request, url, res)
		if err != nil {
			return "", err
		}

		InfoReserve := models.InfoReserve{
			PaymentUrl: models.PaymentUrl{
				Web: res.Info.PaymentUrl.Web,
				App: res.Info.PaymentUrl.App,
			},
			PaymentAccessToken: res.Info.PaymentAccessToken,
			TransactionId:      res.Info.TransactionId,
		}

		data := &models.TransactionLinePay{
			ID:            primitive.NewObjectID(),
			ParkingID:     parking_Id,
			CustomerEmail: email,
			OrderId:       order_id,
			Packages:      pk,
			Info:          InfoReserve,
			Status:        "Pending",
			TimeStamp:     time.Now(),
		}

		_, err = ps.TransactionStorage.InsertLog(ctx, data)
		if err != nil {
			return "", err
		}

		_, err = ps.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{
			"$set": bson.M{
				"transaction_id": res.Info.TransactionId,
			},
		})
		if err != nil {
			return "", err
		}

		return res.Info.PaymentUrl.Web, nil
	}
}

func (ps PaymentServices) LinePayConfirm(ctx context.Context, transactionId, order_id, type_reserve string, is_extend bool) error {
	url := "/v3/payments/" + transactionId + "/confirm"
	log_transaction := new(models.TransactionLinePay)
	err := ps.TransactionStorage.FindLogInterface(ctx, bson.M{"order_id": order_id}, log_transaction)
	if err == mongo.ErrNoDocuments {
		return err
	}

	log_reserve := new(models.Reservation)
	err = ps.ReserveStorage.FindLogInterface(ctx, bson.M{"order_id": order_id}, log_reserve)
	if err == mongo.ErrNoDocuments {
		return err
	}

	request := models.LineConfirmAPIRequest{
		Amount:   int(log_transaction.Packages[0].Amount),
		Currency: "THB",
	}

	res := new(models.LineConfirmAPIResponse)
	err = ps.HttpClient.SendLinePayConfirmAPI(request, url, transactionId, res)
	if err != nil {
		return err
	}

	_, err = ps.TransactionStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"status": "Successful"}})
	if err != nil {
		return err
	}

	if type_reserve == "in_advance" {
		var status string
		if is_extend {
			status = "Parking"
		} else {
			status = "Pending Approval Process"
		}
		_, err = ps.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"status": status}})
		if err != nil {
			return err
		}

		go ps.HttpClient.SendCancelReserveInAdvanceAPI(order_id, log_reserve.CustomerEmail)

	} else {
		var status string
		if is_extend {
			status = "Parking"
		} else {
			status = "Process"
		}
		_, err = ps.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"status": status}})
		if err != nil {
			return err
		}

		go ps.HttpClient.SendTimeoutReserveAPI(order_id, log_reserve.DateEnd, strconv.Itoa(log_reserve.HourEnd), strconv.Itoa(log_reserve.MinEnd))
	}

	err = ps.CustomerStorage.AddPoint(ctx, log_reserve.CustomerEmail, log_reserve.Price)
	if err != nil {
		return err
	}

	_, err = ps.ParkingAreaStorage.UpdateLogByInterface(ctx, bson.M{
		"_id": log_reserve.ParkingID,
	}, bson.M{"$push": bson.M{
		"reserve_log": bson.M{
			"customer_email": log_reserve.CustomerEmail,
			"hour_start":     log_reserve.HourStart,
			"hour_end":       log_reserve.HourEnd,
			"min_start":      log_reserve.MinStart,
			"min_end":        log_reserve.MinEnd,
			"date_start":     log_reserve.DateStart,
			"date_end":       log_reserve.DateEnd,
		},
	}})
	if err != nil {
		return err
	}
	go ps.HttpClient.SendRemoveJobAPI("CR_" + order_id)

	return nil
}

func (ps PaymentServices) LinePayConfirmFine(ctx context.Context, transactionId, order_id string) error {
	url := "/v3/payments/" + transactionId + "/confirm"
	log_transaction := new(models.TransactionLinePay)
	err := ps.TransactionStorage.FindLogInterface(ctx, bson.M{"order_id": order_id}, log_transaction)
	if err == mongo.ErrNoDocuments {
		return err
	}

	err = ps.CustomerStorage.ResetFineCustomer(ctx, log_transaction.CustomerEmail)
	if err != nil {
		return err
	}
	request := models.LineConfirmAPIRequest{
		Amount:   int(log_transaction.Packages[0].Amount),
		Currency: "THB",
	}

	res := new(models.LineConfirmAPIResponse)
	err = ps.HttpClient.SendLinePayConfirmAPI(request, url, transactionId, res)
	if err != nil {
		return err
	}

	_, err = ps.TransactionStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{"$set": bson.M{"status": "Successful"}})
	if err != nil {
		return err
	}

	return nil
}

func (ps PaymentServices) LinePayCancel(ctx context.Context, transactionId, order_id string) error {
	return nil
}

func (ps PaymentServices) CashbackReserve(ctx context.Context, email, order_id, parking_id, type_reserve, action string, req *models.LineReserveRequest) error {
	if action == "fine" {
		log_reserve := new(models.Reservation)
		err := ps.ReserveStorage.FindLogInterface(ctx, bson.M{"order_id": order_id}, log_reserve)
		if err == mongo.ErrNoDocuments {
			return err
		}

		err = ps.CustomerStorage.ResetFineCustomer(ctx, log_reserve.CustomerEmail)
		if err != nil {
			return err
		}
	} else {
		provider_id, err := primitive.ObjectIDFromHex(req.ProviderID)
		if err != nil {
			fmt.Println("Error parking_id not incorrect format:", err)
		}
		parking_Id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error parking_id not incorrect format:", err)
		}
		pk := []models.Package{
			{
				PackageID: req.ProviderID,
				Amount:    req.Quantity * float32(req.Price),
				Name:      req.ParkingName,
				Products: []models.Product{
					{
						Name:     req.ParkingName,
						Quantity: req.Quantity,
						Price:    req.Price,
						ImageURL: "",
					},
				},
			},
		}
		data := &models.TransactionLinePay{
			ID:            primitive.NewObjectID(),
			ParkingID:     parking_Id,
			CustomerEmail: email,
			OrderId:       order_id,
			Packages:      pk,
			Info:          models.InfoReserve{},
			Status:        "Successful",
			TimeStamp:     time.Now(),
		}

		_, err = ps.TransactionStorage.InsertLog(ctx, data)
		if err != nil {
			return err
		}

		log_reserve := new(models.Reservation)
		err = ps.ReserveStorage.FindLogInterface(ctx, bson.M{"order_id": order_id}, log_reserve)
		if err == mongo.ErrNoDocuments {
			return err
		}

		_, err = ps.ParkingAreaStorage.UpdateLogByInterface(ctx, bson.M{
			"_id": log_reserve.ParkingID,
		}, bson.M{"$push": bson.M{
			"reserve_log": bson.M{
				"customer_email": log_reserve.CustomerEmail,
				"hour_start":     log_reserve.HourStart,
				"hour_end":       log_reserve.HourEnd,
				"min_start":      log_reserve.MinStart,
				"min_end":        log_reserve.MinEnd,
				"date_start":     log_reserve.DateStart,
				"date_end":       log_reserve.DateEnd,
			},
		}})
		if err != nil {
			return err
		}
		if type_reserve == "current" {
			var status string
			if log_reserve.IsExtend {
				status = "Parking"
			} else {
				status = "Process"
			}
			_, err = ps.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{
				"$set": bson.M{
					"transaction_id": utility.GenerateTransactionID(),
					"status":         status,
				},
			})
			if err != nil {
				return err
			}
			go ps.HttpClient.SendTimeoutReserveAPI(order_id, log_reserve.DateEnd, strconv.Itoa(log_reserve.HourEnd), strconv.Itoa(log_reserve.MinEnd))

		} else if type_reserve == "in_advance" {
			var status string
			if log_reserve.IsExtend {
				status = "Parking"
			} else {
				status = "Pending Approval Process"
				go ps.NotificationService.ConfirmReserveInAdvanceNotification(ctx, provider_id, email, *log_reserve)
				go ps.HttpClient.SendCancelReserveInAdvanceAPI(order_id, email)
			}
			_, err = ps.ReserveStorage.UpdateLogByInterface(ctx, bson.M{"order_id": order_id}, bson.M{
				"$set": bson.M{
					"transaction_id": utility.GenerateTransactionID(),
					"status":         status,
				},
			})
			if err != nil {
				return err
			}
		}

		err = ps.CustomerStorage.AddPoint(ctx, log_reserve.CustomerEmail, log_reserve.Price)
		if err != nil {
			return err
		}
		if log_reserve.IsExtend {
			go ps.HttpClient.SendRemoveJobAPI("CER_" + order_id)
		} else {
			go ps.HttpClient.SendRemoveJobAPI("CR_" + order_id)

		}
	}

	return nil
}
