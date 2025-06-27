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

func New(gaugeStorage server_storage.GaugeStorage, countStorage server_storage.CounterStorage, mTypes mdata.AvailableMetricsTypes) map[string]http.Handler {
	return Handlers{
		GetListRoute: &GetList{
			gaugeStorage: gaugeStorage,
			countStorage: countStorage,
		},
		UpdateByURLParams: &UpdateHandler{
			AvailableMetricsTypes: mTypes,
			gaugeStorage:          gaugeStorage,
			countStorage:          countStorage,
		},
		GetByURLParams: &GetHandler{
			AvailableMetricsTypes: mTypes,
			gaugeStorage:          gaugeStorage,
			countStorage:          countStorage,
		},
		UpdateByJSON: &JSONUpdateHandler{
			gaugeStorage,
			countStorage,
		},
		GetByJSON: &GetJSONMetricsHandler{
			gaugeStorage: gaugeStorage,
			countStorage: countStorage,
		},
	}
}
