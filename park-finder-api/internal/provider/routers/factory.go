package routers

import (
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	cs "gitlab.com/parking-finder/parking-finder-api/internal/customer/services"
	"gitlab.com/parking-finder/parking-finder-api/internal/httpclient"
	ms "gitlab.com/parking-finder/parking-finder-api/internal/message"
	ns "gitlab.com/parking-finder/parking-finder-api/internal/notification"
	"gitlab.com/parking-finder/parking-finder-api/internal/provider/services"
	rss "gitlab.com/parking-finder/parking-finder-api/internal/reserve"
)

type ProviderRouters struct {
	CustomerServices    cs.ICustomerServices
	ProviderServices    services.IProviderServices
	ReserveService      rss.IReserveServices
	NotificationService ns.INotificationServices
	MessageServices     ms.IMessageServices

	HttpClient *httpclient.HTTPClient
}

func NewProviderRouter(g *echo.Group, ps services.IProviderServices, ns ns.INotificationServices, rss rss.IReserveServices, cs cs.ICustomerServices, ms ms.IMessageServices, ht *httpclient.HTTPClient,
) {
	pr := ProviderRouters{
		ProviderServices:    ps,
		NotificationService: ns,
		ReserveService:      rss,
		CustomerServices:    cs,
		MessageServices:     ms,
		HttpClient:          ht,
	}
	jwtConfig := echojwt.Config{
		ErrorHandler:   pr.errorHandler,
		TokenLookup:    "header:Authorization:Bearer ",
		ParseTokenFunc: pr.parseToken,
	}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)

	//Authenication Service
	g.POST("/register", pr.providerRegisterHandler)                       //doc //test
	g.POST("/login", pr.providerLoginHandler)                             //doc //test
	g.POST("/verify_otp", pr.verifyProviderOTP)                           //doc //test
	g.POST("/logout", pr.providerLogoutHandler, jwtMiddleware)            //doc //test
	g.PATCH("/change_password", pr.updateProviderPassword, jwtMiddleware) //doc //test
	g.POST("/send_forgot_otp", pr.sendOTPForgotHandler)                   //doc //test
	g.POST("/verify_otp_forgot", pr.verifyCustomerOTPForgotHandler)       //doc //test
	g.PATCH("/new_password", pr.updateCustomerPasswordByForgotHandler)    //doc //test
	g.GET("/verify_email/:email", pr.verifyProviderEmail)                 //Doc //test
	g.POST("/resend_register", pr.resendVerifsyHandler)                   //Doc //test
	g.POST("/resend_forgot_otp", pr.resendOTPForgotHandler)               //Doc //test
	g.POST("/resend_login_otp", pr.resendOTPLoginHandler)                 //Doc //test

	//Account Service
	g.GET("/profile", pr.getProviderProfileHandler, jwtMiddleware)       //Doc //test
	g.PATCH("/profile", pr.updateProviderProfileHandler, jwtMiddleware)  //Doc //test
	g.PATCH("/bank_account", pr.updateBankAccountHandler, jwtMiddleware) //Doc //test

	g.GET("/parking_detail", pr.getParkingDetailHandler, jwtMiddleware) //Doc //test
	g.POST("/my_reserve", pr.myReserveParkingHandler, jwtMiddleware)    //Doc //test

	//Parking Area Service
	g.GET("/my_area", pr.getProviderAreaHandler, jwtMiddleware)                              //Doc //test
	g.POST("/register_area_location", pr.providerRegisterAreaLocationHandler, jwtMiddleware) //Doc //test
	g.POST("/register_area_document", pr.providerRegisterAreaDocumentHandler, jwtMiddleware) //Doc //test
	g.PATCH("/area_quick_open_status", pr.updateOpenAreaQuickStatusHandler, jwtMiddleware)   //Doc //test
	g.PATCH("/area_daily_open_status", pr.updateOpenAreaDailyStatusHandler, jwtMiddleware)   //Doc //test
	g.PATCH("/update_price", pr.updateParkingAreaPriceHandler, jwtMiddleware)                //Doc //test

	//Log Service
	g.GET("/retrive_message_list", pr.providerChatListHandler, jwtMiddleware)       //Doc //test
	g.POST("/retrieve_message_log", pr.providerChatLogHandler, jwtMiddleware)       //Doc //test
	g.POST("/notification_list", pr.providerNotificatiobListHandler, jwtMiddleware) //Doc //test
	g.GET("/profit", pr.getProviderProfitHandler, jwtMiddleware)                    //Doc //test

}
