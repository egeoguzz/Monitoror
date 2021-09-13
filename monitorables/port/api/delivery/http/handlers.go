package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	"github.com/monitoror/monitoror/monitorables/port/api/models"
	notify "github.com/monitoror/monitoror/notify"
)

type PortDelivery struct {
	portUsecase api.Usecase
}

func NewPortDelivery(p api.Usecase) *PortDelivery {
	return &PortDelivery{p}
}

func (h *PortDelivery) GetPort(c echo.Context) error {
	// Bind / check Params
	params := &models.PortParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {

		return err
	}

	tile, err := h.portUsecase.Port(params)

	if err != nil {
		if tile.Type == "PORT" {
			notify.ParaMatcher(params.Hostname, true, err.Error())
		}
		return err
	}
	if tile.Type == "PORT" {
		if tile.Status != coreModels.SuccessStatus {

			notify.ParaMatcher(tile.Label, true, string(tile.Status))
		} else {
			notify.ParaMatcher(tile.Label, false, "")
		}
	}
	return c.JSON(http.StatusOK, tile)
}
