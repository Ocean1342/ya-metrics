package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ya-metrics/pkg/mdata"
)

func (h *Handler) UpdateByJSON(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var mReq mdata.Metrics
	err := json.NewDecoder(req.Body).Decode(&mReq)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Wrong request body: %s", err), http.StatusBadRequest)
		return
	}

	err = h.saveJSONData(&mReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Could not save data: %s", err), http.StatusBadRequest)
		return
	}
	updVal, err := h.getUpdatedData(&mReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get updated data: %s", err), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(updVal)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create response body: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func (h *Handler) saveJSONData(mReq *mdata.Metrics) error {
	switch mReq.MType {
	case mdata.COUNTER:
		if mReq.Delta == nil {
			return fmt.Errorf("еmpty %s value", mReq.MType)
		}
		err := h.countStorage.Set(mdata.NewSimpleCounter(mReq.ID, *mReq.Delta))
		fmt.Println("Received: Counter", mReq.ID, *mReq.Delta)
		if err != nil {
			return fmt.Errorf("could not save data in storage")
		}
	case mdata.GAUGE:
		if mReq.Value == nil {
			return fmt.Errorf("еmpty %s value", mReq.MType)
		}
		err := h.gaugeStorage.Set(mdata.NewSimpleGauge(mReq.ID, *mReq.Value))
		fmt.Println("Received: Gauge", mReq.ID, *mReq.Value)
		if err != nil {
			return fmt.Errorf("could not save %s data in storage", mdata.GAUGE)
		}
	default:
		return fmt.Errorf("undefined metric type %s", mReq.MType)
	}

	return nil
}

func (h *Handler) getUpdatedData(mReq *mdata.Metrics) (*mdata.Metrics, error) {
	switch mReq.MType {
	case mdata.COUNTER:
		v, err := h.countStorage.Get(mReq.ID)
		if err != nil {
			return nil, err
		}
		delta := v.GetValue()
		return &mdata.Metrics{
			ID:    v.GetName(),
			MType: mdata.COUNTER,
			Delta: &delta,
		}, nil
	case mdata.GAUGE:
		v := h.gaugeStorage.Get(mReq.ID)
		if v == nil {
			return nil, fmt.Errorf("not found delta with id: %s", mReq.ID)
		}
		value := v.GetValue()
		return &mdata.Metrics{
			ID:    v.GetName(),
			MType: mdata.GAUGE,
			Value: &value,
		}, nil

	default:
		return nil, fmt.Errorf("undefined metric type %s", mReq.MType)
	}
}
