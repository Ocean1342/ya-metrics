package handlers

import (
	"database/sql"
	"go.uber.org/zap"
	"net/http"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type Handlers map[string]http.Handler

type Handler struct {
	availableMetricsTypes mdata.AvailableMetricsTypes
	gaugeStorage          server_storage.GaugeStorage
	countStorage          server_storage.CounterStorage
	db                    *sql.DB
	log                   *zap.SugaredLogger
}

func New(
	gaugeStorage server_storage.GaugeStorage,
	countStorage server_storage.CounterStorage,
	mTypes mdata.AvailableMetricsTypes,
	db *sql.DB,
	log *zap.SugaredLogger,
) *Handler {
	return &Handler{
		availableMetricsTypes: mTypes,
		gaugeStorage:          gaugeStorage,
		countStorage:          countStorage,
		db:                    db,
		log:                   log,
	}
}
