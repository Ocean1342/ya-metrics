// http handlers package
package handlers

import (
	"database/sql"
	"go.uber.org/zap"
	"net/http"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type Handlers map[string]http.Handler

// Handler http handle struct
type Handler struct {
	// availableMetricsTypes - available metrics types e.g. gauge, counter
	availableMetricsTypes mdata.AvailableMetricsTypes
	// gaugeStorage - gauge type storage
	gaugeStorage server_storage.GaugeStorage
	//countStorage - count type storage
	countStorage server_storage.CounterStorage
	// db - database connection
	db *sql.DB
	// log - logger
	log *zap.SugaredLogger
}

// New construct Handler
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
