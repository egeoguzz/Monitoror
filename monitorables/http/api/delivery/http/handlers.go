package http

import (
	netHttp "net/http"

	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"
	notify "github.com/monitoror/monitoror/notify"
)

//nolint:golint
type HTTPDelivery struct {
	httpUsecase api.Usecase
}

func NewHTTPDelivery(p api.Usecase) *HTTPDelivery {
	return &HTTPDelivery{p}
}

func (h *HTTPDelivery) GetHTTPStatus(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPStatusParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {

		return err
	}

	tile, err := h.httpUsecase.HTTPStatus(params)

	if err != nil {
		if tile.Type == "HTTP-STATUS" {
			notify.ParaMatcher(params.URL, true, err.Error())

		}

		//	notify.ParaMatcher(tile.Label, string(tile.Type), false)
		return err
	}
	if tile.Type == "HTTP-STATUS" {
		if tile.Status != coreModels.SuccessStatus {
			notify.ParaMatcher(tile.Label, true, string(tile.Status))
		} else {
			notify.ParaMatcher(tile.Label, false, "")
		}
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPRaw(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPRawParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.httpUsecase.HTTPRaw(params)
	if err != nil {
		if tile.Type == "HTTP-RAW" {
			notify.ParaMatcher(params.URL, true, err.Error())
		}

		return err
	}
	if tile.Type == "HTTP-RAW" {
		if tile.Status != coreModels.SuccessStatus {
			notify.ParaMatcher(tile.Label, true, string(tile.Status))
		} else {
			notify.ParaMatcher(tile.Label, false, "")
		}
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPFormatted(c echo.Context) error {
	// Bind / Check Params
	params := &models.HTTPFormattedParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.httpUsecase.HTTPFormatted(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}
