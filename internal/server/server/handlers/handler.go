package handlers

import (
	"net/http"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type Handlers map[string]http.Handler

const GetListRoute = "getList"
const UpdateByURLParams = "updateByURLParams"
const GetByURLParams = "getByURLParams"
const UpdateByJSON = "updateByJSON"
const GetByJSON = "getByJSON"

type Handler struct {
	availableMetricsTypes mdata.AvailableMetricsTypes
	gaugeStorage          server_storage.GaugeStorage
	countStorage          server_storage.CounterStorage
}

func New(gaugeStorage server_storage.GaugeStorage, countStorage server_storage.CounterStorage, mTypes mdata.AvailableMetricsTypes) *Handler {
	return &Handler{
		availableMetricsTypes: mTypes,
		gaugeStorage:          gaugeStorage,
		countStorage:          countStorage,
	}
}
