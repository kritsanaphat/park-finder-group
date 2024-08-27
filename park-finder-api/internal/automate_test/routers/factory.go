package routers

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/internal/automate_test/services"
)

type AutomateTestRouter struct {
	AutomateTeserServices services.IAutomateTeserServices
}

func NewAutomateTestRouter(g *echo.Group, ats services.IAutomateTeserServices) {
	atr := AutomateTestRouter{
		AutomateTeserServices: ats,
	}

	//Customer Authenication Service
	g.POST("/reset", atr.resetData)
	g.POST("/add_cashback", atr.addCashback) //Doc //test

}
