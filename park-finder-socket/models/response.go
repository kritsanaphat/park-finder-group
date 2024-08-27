package models

import (
	"github.com/labstack/echo/v4"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token" validate:"required"`
	Data        *echo.Map `json:"data,omitempty"`
}

type CustomerProfileResponse struct {
	Profile *echo.Map `json:"profile"`
	Status  int       `json:"status"`
}

type CustomerAddressResponse struct {
	Data   []CustomerAddress `json:"address"`
	Status int               `json:"status"`
}

type CustomerCarResponse struct {
	Data   []Car `json:"data"`
	Status int   `json:"status"`
}

type SearchAreaResponse struct {
	Data   *[]ParkingArea `json:"data"`
	Status int            `json:"status"`
}

type ParkingDetailResponse struct {
	Data   *ParkingArea `json:"data"`
	Status int          `json:"status"`
}

type LineReserveAPIResponse struct {
	ReturnCode    string      `json:"returnCode"`
	ReturnMessage string      `json:"returnMessage"`
	Info          InfoReserve `json:"info"`
}

type LineConfirmAPIResponse struct {
	ReturnCode    string      `json:"returnCode"`
	ReturnMessage string      `json:"returnMessage"`
	Info          InfoConfirm `json:"info"`
}

type ReserveResponse struct {
	OrderID string `json:"order_id"`
	Status  int    `json:"status"`
}

type MyReserveResponse struct {
	Reservation []echo.Map `json:"data"`
	Status      int        `json:"status"`
}

type ProfitResponse struct {
	Data   *Profit `json:"data"`
	Status int     `json:"status"`
}
