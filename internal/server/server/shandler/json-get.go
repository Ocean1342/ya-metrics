package shandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type GetJSONMetricsHandler struct {
	gaugeStorage server_storage.GaugeStorage
	countStorage server_storage.CounterStorage
}

func (g *GetJSONMetricsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var getReq mdata.Metrics
	err := json.NewDecoder(req.Body).Decode(&getReq)
	if err != nil {
		http.Error(w, "Could not decode request body", http.StatusMethodNotAllowed)
		return
	}
	m, err := g.getMetric(&getReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("id: %s, mtype: %s not found", getReq.ID, getReq.MType), http.StatusNotFound)
		return
	}
	data, err := json.Marshal(m)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create response body: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (g *GetJSONMetricsHandler) getMetric(getReq *mdata.Metrics) (*mdata.Metrics, error) {
	switch getReq.MType {
	case mdata.GAUGE:
		v := g.gaugeStorage.Get(getReq.ID)
		if v == nil {
			return nil, fmt.Errorf("id: %s, mtype: %s not found", getReq.ID, getReq.MType)
		}
		value := v.GetValue()
		return &mdata.Metrics{
			ID:    getReq.ID,
			MType: getReq.MType,
			Value: &value,
		}, nil
	case mdata.COUNTER:
		v, err := g.countStorage.Get(getReq.ID)
		if err != nil || v == nil {
			return nil, fmt.Errorf("id: %s, mtype: %s not found", getReq.ID, getReq.MType)
		}
		value := v.GetValue()
		return &mdata.Metrics{
			ID:    getReq.ID,
			MType: getReq.MType,
			Delta: &value,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported mtype: %s", getReq.MType)
	}
}
