package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
)

func (ar AutomateTestRouter) resetData(c echo.Context) error {

	err := ar.AutomateTeserServices.ResetData()
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "ok"})
}

func (ar AutomateTestRouter) addCashback(c echo.Context) error {

	err := ar.AutomateTeserServices.AddCashback()
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "ok"})
}
