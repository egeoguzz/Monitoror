package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	"github.com/monitoror/monitoror/monitorables/ping/api/models"
	notify "github.com/monitoror/monitoror/notify"
)

type PingDelivery struct {
	pingUsecase api.Usecase
}

func NewPingDelivery(p api.Usecase) *PingDelivery {
	return &PingDelivery{p}
}

func (h *PingDelivery) GetPing(c echo.Context) error {
	// Bind / Check Params
	params := &models.PingParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.pingUsecase.Ping(params)
	if err != nil {
		if tile.Type == "PING" {
			notify.ParaMatcher(params.Hostname, true, err.Error())
		}

		return err
	}
	if tile.Type == "PING" {
		if tile.Status != coreModels.SuccessStatus {
			notify.ParaMatcher(tile.Label, true, string(tile.Status))
		} else {
			notify.ParaMatcher(tile.Label, false, "")
		}
	}
	return c.JSON(http.StatusOK, tile)
}
