package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// API
func (ns NotificationServices) NotificationList(ctx context.Context, receiver_id, type_client string) ([]models.Notification, error) {
	notification_list, err := ns.NotificationStorage.FindNotification(ctx, receiver_id, type_client)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(notification_list)-1; i < j; i, j = i+1, j-1 {
		notification_list[i], notification_list[j] = notification_list[j], notification_list[i]
	}
	return notification_list, nil
}

// Customer
func (ns NotificationServices) ProviderConfirmReserveInAdvanceNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation, address string) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType:  "Personal",
		ReceiverID:     &receiver_id,
		Title:          "การจองของคุณได้รับการอนุมัติ",
		Description:    fmt.Sprintf("%s|%s|%s|%s|%d|%d|%d|%d|cashback=%d|%s", address, reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, ""),
		CallbackMethod: []models.CallbackMethod{},
	}
	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Customer
func (ns NotificationServices) ProviderCancelReserveInAdvanceNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation, address string) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType:  "Personal",
		ReceiverID:     &receiver_id,
		Title:          "การจองของคุณถูกปฏิเสธ",
		Description:    fmt.Sprintf("%s|%s|%s|%s|%d|%d|%d|%d|cashback=%d|%s", address, reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, ""),
		CallbackMethod: []models.CallbackMethod{},
	}
	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Customer
func (ns NotificationServices) BeforeTimeOutReserveNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType: "Personal",
		ReceiverID:    &receiver_id,
		Title:         "เวลาการจองของคุณกำลังจะหมด",
		Description:   fmt.Sprintf("คุณต้องการที่จะขยายเวลา เพิ่ม 1 ชั่วโมงหรือไม่?|%s|%s|%s|%d|%d|%d|%d|ราคาที่ต้องชำระเพิ่ม=%d|%s", reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, ""),
		CallbackMethod: []models.CallbackMethod{

			{
				Action:      "Extend",
				CallBackURL: fmt.Sprintf("http://%s/customer/extend_reserve?order_id=%s&action=normal", os.Getenv("HOST"), reserve.OrderID),
			},
		},
	}

	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Customer
func (ns NotificationServices) TimeOutReserveNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType: "Personal",
		ReceiverID:    &receiver_id,
		Title:         "การจองของคุณครบกำหนด",
		Description:   fmt.Sprintf("คุณต้องการทำการจอดต่อหรือไม่?|%s|%s|%s|%d|%d|%d|%d|ราคาที่ต้องชำระเพิ่ม=%d|%s", reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, ""),
		CallbackMethod: []models.CallbackMethod{
			{
				Action:      "Extend",
				CallBackURL: fmt.Sprintf("http://%s/customer/extend_reserve?order_id=%s&action=normal", os.Getenv("HOST"), reserve.OrderID),
			},
		},
	}

	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Customer
func (ns NotificationServices) AfterTimeOutReserveNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation, fine int) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType: "Personal",
		ReceiverID:    &receiver_id,
		Title:         "คุณจอดเกินเวลาที่ได้ทำการจองมา",
		Description:   fmt.Sprintf("กรุณาจ่ายค่าปรับที่ค้างอยู่?|%s|%s|%s|%d|%d|%d|%d|ค่าปรับเป็นจำนวนเงิน=%d|%s", reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.MinStart, reserve.MinEnd, reserve.HourEnd, fine, ""),
		CallbackMethod: []models.CallbackMethod{
			{
				Action:      "Pay",
				CallBackURL: fmt.Sprintf("http://%s/customer/extend_reserve?order_id=%s&action=normal", os.Getenv("HOST"), reserve.OrderID),
			},
		},
	}

	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Customer
func (ns NotificationServices) ReservationCancelNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType:  "Personal",
		ReceiverID:     &receiver_id,
		Title:          "ที่จอดที่คุณจองมาเกิดปัญหาบางอย่าง",
		Description:    fmt.Sprintf("ทางเราต้องขออภัยกับเหตุการณ์ที่เกิดขึ้น|%s|%s|%s|%d|%d|%d|%d|cashback=%d|%s", reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, ""),
		CallbackMethod: []models.CallbackMethod{},
	}

	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Customer
func (ns NotificationServices) LeaveTimeOutReserveNotification(ctx context.Context, receiver_id, parking_id primitive.ObjectID, order_id, parking_name string) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType: "Personal",
		ReceiverID:    &receiver_id,
		Title:         "การจองของคุณสำเร็จแล้ว",
		Description:   fmt.Sprintf("รถของคุณได้ออกจากที่จอดเรียบร้อยแล้ว โปรด review การจองของคุณ ภายใน 1 ชั่วโมง|%s|%s", order_id, parking_id.Hex()),
		CallbackMethod: []models.CallbackMethod{
			{
				Action:      "Review",
				CallBackURL: parking_id.Hex(),
			},
		},
	}

	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Customer
func (ns NotificationServices) VertifyCustomerCarNotification(ctx context.Context, receiver_id primitive.ObjectID, license_plate, module_code, pic_url string, reserve *models.Reservation) error {
	fmt.Println("hi")
	fmt.Println("--------Notification-------- :", receiver_id.Hex())
	notification := models.Notification{
		BroadcastType: "Personal",
		ReceiverID:    &receiver_id,
		Title:         "เราไม่พบรถของคุณ",
		Description:   fmt.Sprintf("กรุณายืนยันว่าใช่รถของคุณหรือไม่?|%s|%s|%s|%d|%d|%d|%d|%d|%s", reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, pic_url),
		CallbackMethod: []models.CallbackMethod{
			{
				Action:      "Confirm",
				CallBackURL: fmt.Sprintf("http://%s/customer/start_reserve?parking_id=%s&license_plate=%s&module_code=%s&action=force", os.Getenv("HOST"), reserve.ParkingID.Hex(), license_plate, module_code),
			},
			{
				Action: "Not me",
				CallBackURL: fmt.Sprintf("http://%s/customer/report_verify?order_id=%s&provider_id=%s",
					os.Getenv("HOST"),
					reserve.OrderID,
					receiver_id.Hex()),
			},
		},
	}
	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Provider
func (ns NotificationServices) ConfirmReserveInAdvanceNotification(ctx context.Context, receiver_id primitive.ObjectID, customer_email string, reserve models.Reservation) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType: "Personal",
		ReceiverID:    &receiver_id,
		Title:         "โปรดยืนยันการจองของผู้ใช้งาน",
		Description:   fmt.Sprintf("โปรดตอบกลับภายใน 24 ชั่วโมง|%s|%s|%s|%d|%d|%d|%d|%d|%s", reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, ""),
		CallbackMethod: []models.CallbackMethod{
			{
				Action:      "Confirm",
				CallBackURL: fmt.Sprintf("http://%s/webhook/internal/notification/confirm_reserve_in_advance_notification/confirm?order_id=%s&email=%s", os.Getenv("HOST"), reserve.OrderID, customer_email),
			},
			{
				Action:      "Cancel",
				CallBackURL: fmt.Sprintf("http://%s/webhook/internal/notification/confirm_reserve_in_advance_notification/cancel?order_id=%s&email=%s", os.Getenv("HOST"), reserve.OrderID, customer_email),
			},
		},
	}
	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Provider
func (ns NotificationServices) ParkingAreaStatusUpdateNotification(ctx context.Context, receiver_id primitive.ObjectID, status, parking_name string) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType:  "Personal",
		ReceiverID:     &receiver_id,
		Title:          "ที่จอดรถของคุณมีการอัพเดทสถานะ",
		Description:    fmt.Sprintf("ที่จอดรถ %s ได้รับการอัพเดทสถานนะเป็น  %s", parking_name, status),
		CallbackMethod: []models.CallbackMethod{},
	}
	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Provider
func (ns NotificationServices) ReportParkingAreaNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error {
	fmt.Println("--------Notification-------- :", receiver_id.Hex())

	notification := models.Notification{
		BroadcastType:  "Personal",
		ReceiverID:     &reserve.ProviderID,
		Title:          "เกิดปํญหาในที่จอดรถของคุณ",
		Description:    fmt.Sprintf("เกิดปัญหาเนื่องจากมีผู้อื่นมาจอดรถในพื้่นที่ของคุณ โปรดตรวจสอบความถูกต้อง|%s|%s|%s|%d|%d|%d|%d|%d|%s", reserve.ParkingName, reserve.DateStart, reserve.DateEnd, reserve.HourStart, reserve.HourEnd, reserve.MinStart, reserve.MinEnd, reserve.Price, ""),
		CallbackMethod: []models.CallbackMethod{},
	}
	fmt.Println("Sending Notification To :", receiver_id.Hex())

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}

// Broadcast
func (ns NotificationServices) AddRewardNotification(ctx context.Context, tile, description, id, preview_url string) error {
	fmt.Println("--------Notification-------- :", "Customer")

	notification := models.Notification{
		BroadcastType: "customer",
		ReceiverID:    nil,
		Title:         tile,
		Description:   description,
		CallbackMethod: []models.CallbackMethod{
			{
				Action:      "More",
				CallBackURL: fmt.Sprintf("http://%s/customer/reward_detail?_id=%s", os.Getenv("HOST"), id),
			},
		},
	}
	fmt.Println("Sending Notification To : Customer")

	byteValue, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	byteKey, err := json.Marshal("notification")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err := ns.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
		fmt.Println("Error producing message:", err)
		return err
	}
	return nil
}
