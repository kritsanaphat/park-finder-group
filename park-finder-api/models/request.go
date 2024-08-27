package models

type RegisterAccountRequest struct {
	FristName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Phone     string `json:"phone" bson:"phone"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

type AddRewardRequest struct {
	Name            string   `json:"name"`
	Point           int      `json:"point"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	ExpiredDate     string   `json:"expired_date"`
	PreviewImageURL string   `json:"preview_url"`
	Webhook         string   `json:"webhook"`
	Condition       []string `json:"condition"`
	QuotaCount      int      `json:"quota_count"`
}

type RegisterCarRequest struct {
	Name          string `json:"name" bson:"name"`
	LicensePlate  string `json:"license_plate" bson:"license_plate"`
	Brand         string `json:"brand" bson:"brand"`
	Model         string `json:"model" bson:"model"`
	Color         string `json:"color" bson:"color"`
	CarPictureURL string `json:"car_picture_url" bson:"car_picture_url"`
}

type RegisterParkingAreaFirstStepRequest struct {
	ParkingName  string   `json:"parking_name" bson:"parking_name"`
	AddressText  string   `json:"address_text" bson:"address_text"`
	Sub_district string   `json:"sub_district" bson:"sub_district"`
	District     string   `json:"district" bson:"district"`
	Province     string   `json:"province" bson:"province"`
	Postal_code  string   `json:"postal_code" bson:"postal_code"`
	Tag          []string `json:"tag" bson:"tag"`
	Latitude     float64  `json:"latitude" bson:"latitude"`
	Longitude    float64  `json:"longitude" bson:"longitude"`
}

type RegisterParkingAreaDocumentStepRequest struct {
	ParkingPictureUrl     string   `json:"parking_picture_url" bson:"parking_picture_url"`
	TitleDeedUrl          string   `json:"title_deed_url" bson:"title_deed_url"`
	LandCertificateUrl    string   `json:"land_certificate_url" bson:"land_certificate_url"`
	IDCardUrl             string   `json:"id_card_url" bson:"id_card_url"`
	ToatalParkingCount    int      `json:"total_parking_count" bson:"total_parking_count"`
	OverviewPictureUrl    []string `json:"over_view_picture_url" bson:"over_view_picture_url"`
	MeasurementPictureUrl []string `json:"measurement_picture_url" bson:"measurement_picture_url"`
	Price                 int16    `json:"price" bson:"price"`
	StatusApply           string   `json:"status_apply" bson:"status_apply"`
}

type RegisterParkingAreaRequest struct {
	ParkingName        string  `json:"parking_name" bson:"parking_name"`
	AddressText        string  `json:"address_text" bson:"address_text"`
	Sub_district       string  `json:"sub_district" bson:"sub_district"`
	District           string  `json:"district" bson:"district"`
	Province           string  `json:"province" bson:"province"`
	Postal_code        string  `json:"postal_code" bson:"postal_code"`
	ParkingPictureUrl  string  `json:"parking_picture_url" bson:"parking_picture_url"`
	TitleDeedUrl       string  `json:"title_deed_url" bson:"title_deed_url"`
	LandCertificateUrl string  `json:"land_certificate_url" bson:"land_certificate_url"`
	IDCardUrl          string  `json:"id_card_url" bson:"id_card_url"`
	ToatalParkingCount int     `json:"total_parking_count" bson:"total_parking_count"`
	Tag                string  `json:"tag" bson:"tag"`
	Price              int16   `json:"price" bson:"price"`
	Latitude           float64 `json:"latitude" bson:"latitude"`
	Longitude          float64 `json:"longitude" bson:"longitude"`
}

type UpdateProfileRequest struct {
	SSN               string `json:"ssn" bson:"ssn"`
	FirstName         string `json:"first_name" bson:"first_name"`
	LastName          string `json:"last_name" bson:"last_name"`
	Birthday          string `json:"birth_day" bson:"birth_day"`
	Phone             string `json:"phone" bson:"phone"`
	ProfilePictureURL string `json:"profile_picture_url" bson:"profile_picture_url"`
}

type UpdateOpenAreaDailyStatusRequest struct {
	Type          string         `json:"type"`
	ParkingAreaID string         `json:"_id"`
	Monday        OpenTimeDetail `json:"monday" bson:"monday"`
	Tuesday       OpenTimeDetail `json:"tuesday" bson:"tuesday"`
	Wednesday     OpenTimeDetail `json:"wednesday" bson:"wednesday"`
	Thursday      OpenTimeDetail `json:"thursday" bson:"thursday"`
	Friday        OpenTimeDetail `json:"friday" bson:"friday"`
	Saturday      OpenTimeDetail `json:"saturday" bson:"saturday"`
	Sunday        OpenTimeDetail `json:"sunday" bson:"sunday"`
}

type UpdateOpenAreaQuickStatusRequest struct {
	ParkingAreaID string `json:"_id"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Range         int    `json:"range"`
}

type UpdateOpenAreaInAdvanceStatusRequest struct {
	ParkingAreaID string   `json:"_id"`
	Type          string   `json:"type"`
	Status        string   `json:"status"`
	Date          []string `json:"date"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordEmailRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type OTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type LocationRequest struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type RegisterAddressRequest struct {
	Address      string  `json:"address" bson:"address"`
	LocationName string  `json:"location_name" bson:"location_name"`
	Latitude     float64 `json:"latitude" bson:"latitude"`
	Longitude    float64 `json:"longitude" bson:"longitude"`
}

type SearchQueryRequest struct {
	Keyword   string   `json:"keyword"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Status    bool     `json:"status"`
	Review    int16    `json:"review"`
	MaxPrice  int16    `json:"max_price"`
	MinPrice  int16    `json:"min_price"`
	Date      []string `json:"date"`
	HourStart int      `json:"hour_start"`
	HourEnd   int      `json:"hour_end"`
	MinStart  int      `json:"min_start"`
	MinEnd    int      `json:"min_end"`
}

type LineReserveRequest struct {
	OrderID     string  `json:"order_id"`
	ParkingID   string  `json:"parking_id"`
	ProviderID  string  `json:"provider_id"`
	Quantity    float32 `json:"quantity"`
	Price       int     `json:"price"`
	ParkingName string  `json:"parking_name"`
	CashBack    int     `json:"cashback"`
}

type LineReserveAPIRequest struct {
	Amount       float32      `json:"amount"`
	Currency     string       `json:"currency"`
	OrderId      string       `json:"orderId"`
	Packages     []Package    `json:"packages"`
	RedirectUrls RedirectUrls `json:"redirectUrls"`
}

type LineConfirmAPIRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type ReserveRequest struct {
	ProviderID    string `json:"provider_id"`
	ParkingID     string `json:"parking_id"`
	CarID         string `json:"car_id"`
	DateStart     string `json:"date_start"`
	DateEnd       string `json:"date_end"`
	Type          string `json:"type"`
	HourStart     int    `json:"hour_start"`
	MinStart      int    `json:"min_start"`
	HourEnd       int    `json:"hour_end" `
	MinEnd        int    `json:"min_end"`
	PaymentChanel string `json:"payment_chanel"`
	Price         int    `json:"price"`
}

type MyReserveRequest struct {
	ParkingID string `json:"parking_id"`
	Status    string `json:"status"`
}

type ProfitRequest struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	ParkingID string `json:"parking_id"`
}

type ReviewRequest struct {
	OrderID     string `json:"order_id"`
	Comment     string `json:"comment"`
	ReviewScore int    `json:"review_score"`
	ParkingID   string `json:"parking_id"`
}

type ReportRequest struct {
	OrderID string `json:"order_id"`
	Content string `json:"content"`
}

type SubmitReceiptRequest struct {
	ReceiptImageUrl string `json:"receipt_image_url"`
	ProviderID      string `json:"provider_id"`
	Price           int    `json:"price"`
	Month           string `json:"month"`
	Year            string `json:"year"`
}
