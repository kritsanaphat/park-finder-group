package routers

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	cs "gitlab.com/parking-finder/parking-finder-api/internal/customer/services"
	ms "gitlab.com/parking-finder/parking-finder-api/internal/message"
	ns "gitlab.com/parking-finder/parking-finder-api/internal/notification"
	pms "gitlab.com/parking-finder/parking-finder-api/internal/payment"
	rss "gitlab.com/parking-finder/parking-finder-api/internal/reserve"
	ss "gitlab.com/parking-finder/parking-finder-api/internal/search"
)

type CustomerRouters struct {
	CustomerServices    cs.ICustomerServices
	SearchServices      ss.ISearchServices
	PaymentService      pms.IPaymentServices
	ReserveService      rss.IReserveServices
	MessageServices     ms.IMessageServices
	NotificationService ns.INotificationServices
}

func NewCustomerRouter(g *echo.Group, cs cs.ICustomerServices, ss ss.ISearchServices, pms pms.IPaymentServices, rss rss.IReserveServices, ms ms.IMessageServices, ns ns.INotificationServices) {
	cr := CustomerRouters{
		CustomerServices:    cs,
		SearchServices:      ss,
		PaymentService:      pms,
		ReserveService:      rss,
		MessageServices:     ms,
		NotificationService: ns,
	}
	jwtConfig := echojwt.Config{
		ErrorHandler:   cr.errorHandler,
		TokenLookup:    "header:Authorization:Bearer ",
		ParseTokenFunc: cr.parseToken,
	}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)
	//Customer Authenication Service
	g.POST("/register", cr.customerRegisterHandler)                       //Doc //test
	g.POST("/login", cr.customerLoginHandler)                             //Doc //test
	g.POST("/verify_otp", cr.verifyCustomerOTP)                           //Doc //test
	g.POST("/logout", cr.customerLogoutHandler, jwtMiddleware)            //Doc //test
	g.PATCH("/change_password", cr.updateCustomerPassword, jwtMiddleware) //Doc //test

	g.POST("/send_forgot_otp", cr.sendOTPForgotHandler)                //Doc //test
	g.POST("/verify_otp_forgot", cr.verifyCustomerOTPForgotHandler)    //Doc //test
	g.PATCH("/new_password", cr.updateCustomerPasswordByForgotHandler) //Doc //test

	g.POST("/resend_register", cr.resendVerifsyHandler)     //Doc //test
	g.POST("/resend_forgot_otp", cr.resendOTPForgotHandler) //Doc //test
	g.POST("/resend_login_otp", cr.resendOTPLoginHandler)   //Doc //test

	g.GET("/verify_email/:email", cr.verifyCustomerEmail) //Doc //test

	//Customer Account Service
	g.PATCH("/profile", cr.updateCustomerProfileHandler, jwtMiddleware) //Doc //test
	g.GET("/profile", cr.customerProfileHandler, jwtMiddleware)         //Doc //test

	g.POST("/car", cr.customerRegisteCarHandler, jwtMiddleware)                //Doc //test
	g.GET("/car", cr.customerCarlHandler, jwtMiddleware)                       //Doc //test
	g.PATCH("/car", cr.updateCustomerCarHandler, jwtMiddleware)                //Doc //test
	g.DELETE("/car", cr.deleteCustomerCarHandler, jwtMiddleware)               //Doc //test
	g.PATCH("/car_default", cr.updateCustomerDefaultCarHandler, jwtMiddleware) //Doc //test

	g.GET("/favorite_area", cr.customerFavoriteAddressHandler, jwtMiddleware)     //Doc //test
	g.PATCH("/favorite_area", cr.updateCustomerFavoriteCarHandler, jwtMiddleware) //Doc //test

	g.GET("/address", cr.customerAddressHandler, jwtMiddleware)                        //Doc //test
	g.POST("/address", cr.customerRegisterAddressHandler, jwtMiddleware)               //Doc //test
	g.PATCH("/address", cr.updateCustomerAddressHandler, jwtMiddleware)                //Doc //test
	g.DELETE("/address", cr.deleteCustomerAddressHandler, jwtMiddleware)               //Doc //test
	g.PATCH("/address_default", cr.updateCustomerDefaultAddressHandler, jwtMiddleware) //Doc //test

	// Search service
	g.POST("/search_parking", cr.searchLocationParkingAreaHandler, jwtMiddleware) //Doc //test
	g.GET("/parking_detail", cr.getParkingDetailHandler, jwtMiddleware)           //Doc //test

	// Reserve service
	g.POST("/reserve", cr.reserveParkingHandler, jwtMiddleware)                   //Doc //test
	g.POST("/start_reserve", cr.startReserveParkingHandler, jwtMiddleware)        //Doc //test
	g.POST("/my_reserve", cr.myReserveParkingHandler, jwtMiddleware)              //Doc //test
	g.POST("/my_reserve_detail", cr.myReserveParkingDetailHandler, jwtMiddleware) //Doc //test
	g.POST("/report_verify", cr.reportParkingVerify, jwtMiddleware)               //Doc //test
	g.POST("/review", cr.createReviewHandler, jwtMiddleware)                      //Doc //test
	g.POST("/check_can_review", cr.checkCanReview, jwtMiddleware)                 //Doc //test
	g.POST("/report_reserve", cr.reportReserve, jwtMiddleware)                    //Doc //test
	g.POST("/check_fine", cr.checkFine, jwtMiddleware)                            //Doc //test
	g.POST("/extend_reserve", cr.extendReserveHandler, jwtMiddleware)             //Doc //test
	g.POST("/capture_picture", cr.captureCarReserve, jwtMiddleware)               //Doc //test

	//Payment service
	g.POST("/line-pay/payment", cr.linePayReserveHandler, jwtMiddleware) //Doc //test

	//Reward service
	g.GET("/reward", cr.customerRewardHandler, jwtMiddleware)            //Doc //test
	g.GET("/reward_detail", cr.customerRewardDetail, jwtMiddleware)      //Doc //test
	g.GET("/redeem_reward", cr.customerRedeemReward, jwtMiddleware)      //Doc //test
	g.GET("/my_redeem_reward", cr.customerMyRedeemReward, jwtMiddleware) //Doc //test
	g.GET("/history_point", cr.MyHistoryPoint, jwtMiddleware)            //Doc //test

	//Log service
	g.POST("/retrieve_message_log", cr.customerChatLogHandler, jwtMiddleware)       //Doc //test
	g.POST("/notification_list", cr.customerNotificatiobListHandler, jwtMiddleware) //Doc //test

}
