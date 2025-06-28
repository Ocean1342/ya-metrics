package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ya-metrics/pkg/mdata"
)

func (h *Handler) GetByJSON(w http.ResponseWriter, req *http.Request) {
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
	m, err := h.getMetric(&getReq)
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

func (h *Handler) getMetric(getReq *mdata.Metrics) (*mdata.Metrics, error) {
	switch getReq.MType {
	case mdata.GAUGE:
		v := h.gaugeStorage.Get(getReq.ID)
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
		v, err := h.countStorage.Get(getReq.ID)
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
