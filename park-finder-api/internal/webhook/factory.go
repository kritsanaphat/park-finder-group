package routers

import (
	"github.com/labstack/echo/v4"
	cs "gitlab.com/parking-finder/parking-finder-api/internal/customer/services"
	"gitlab.com/parking-finder/parking-finder-api/internal/httpclient"
	ns "gitlab.com/parking-finder/parking-finder-api/internal/notification"
	pms "gitlab.com/parking-finder/parking-finder-api/internal/payment"
	rs "gitlab.com/parking-finder/parking-finder-api/internal/reserve"
	"gitlab.com/parking-finder/parking-finder-api/pkg/connector"
)

type WebhookRouters struct {
	CustomerServices    cs.ICustomerServices
	PaymentServices     pms.IPaymentServices
	ReserveServices     rs.IReserveServices
	NotificationService ns.INotificationServices
	HttpClient          *httpclient.HTTPClient
	Redis               *connector.Redis
}

func NewWebhookRouter(g *echo.Group, cs cs.ICustomerServices, pms pms.IPaymentServices, rs rs.IReserveServices, ns ns.INotificationServices, ht *httpclient.HTTPClient, rd *connector.Redis,
) {
	wr := WebhookRouters{
		CustomerServices:    cs,
		PaymentServices:     pms,
		ReserveServices:     rs,
		NotificationService: ns,
		HttpClient:          ht,
		Redis:               rd,
	}

	g.GET("/line-pay/reserve/callback", wr.LinePayConfirmCallBackWebhookHandler)
	g.GET("/line-pay/reserve/callback/fine", wr.LinePayConfirmFineCallBackWebhookHandler)
	g.GET("/line-pay/reserve/callback/cancel", wr.LinePayCancelCallBackWebhookHandler)

	g.POST("/internal/cronjob/before_timeout_reserve", wr.BeforeTimeOutReserve)
	g.POST("/internal/cronjob/timeout_reserve", wr.TimeOutReserve)
	g.POST("/internal/cronjob/after_timeout_reserve", wr.AfterTimeOutReserve)

	g.GET("/internal/notification/confirm_reserve_in_advance_notification/confirm", wr.ConfirmReserveInAdvance)
	g.GET("/internal/notification/confirm_reserve_in_advance_notification/cancel", wr.CancelReserveInAdvance)

}
