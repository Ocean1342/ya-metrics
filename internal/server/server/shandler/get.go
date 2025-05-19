package shandler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

func NewGetHandler(
	amt mdata.AvailableMetricsTypes,
	gaugeStorage server_storage.GaugeStorage,
	countStorage server_storage.CounterStorage,
) *GetHandler {
	return &GetHandler{
		AvailableMetricsTypes: amt,
		gaugeStorage:          gaugeStorage,
		countStorage:          countStorage,
	}
}

type GetHandler struct {
	AvailableMetricsTypes mdata.AvailableMetricsTypes
	gaugeStorage          server_storage.GaugeStorage
	countStorage          server_storage.CounterStorage
}

func (gh *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	typeName := chi.URLParam(r, "type")

	if !gh.AvailableMetricsTypes.Isset(typeName) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := chi.URLParam(r, "name")
	g := gh.gaugeStorage.Get(name)
	if g != nil {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf("%v", g.GetValue())))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	c, err := gh.countStorage.Get(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if c != nil {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf("%v", c.GetValue())))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(""))
}
