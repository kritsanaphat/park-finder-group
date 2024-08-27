package routers

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/internal/admin/services"
	ns "gitlab.com/parking-finder/parking-finder-api/internal/notification"
)

type AdminRouters struct {
	AdminServices       services.IAdminServices
	NotificationService ns.INotificationServices
}

func NewAdminRouter(g *echo.Group, as services.IAdminServices, ns ns.INotificationServices) {
	ar := AdminRouters{
		AdminServices:       as,
		NotificationService: ns,
	}
	jwtConfig := echojwt.Config{
		ErrorHandler:   ar.errorHandler,
		TokenLookup:    "header:Authorization:Bearer ",
		ParseTokenFunc: ar.parseToken,
	}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)

	g.POST("/register", ar.adminRegisterHandler)            //Doc //test
	g.POST("/login", ar.adminLoginHandler)                  //Doc //test
	g.POST("/logout", ar.adminLogoutHandler, jwtMiddleware) //Doc //test

	g.POST("/reward", ar.addRewardHandler, jwtMiddleware) //Doc //test

	g.GET("/get_parking_area", ar.adminGetAreaHandler, jwtMiddleware)                //Doc //test
	g.PATCH("/update_patking_area_status", ar.adminUpdateAreaHandler, jwtMiddleware) //Doc //test

	g.GET("/transaction_list", ar.adminGetTransaction, jwtMiddleware) //Doc //test
	g.POST("/submit_receipt", ar.adminSubmitReceipt, jwtMiddleware)   //Doc //test

}
