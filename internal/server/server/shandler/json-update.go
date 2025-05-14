package shandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type JsonUpdateHandler struct {
	gaugeStorage server_storage.GaugeStorage
	countStorage server_storage.CounterStorage
}

func NewJsonUpdateHandler(
	gaugeStorage server_storage.GaugeStorage,
	countStorage server_storage.CounterStorage,
) *JsonUpdateHandler {
	return &JsonUpdateHandler{
		gaugeStorage,
		countStorage,
	}
}

func (j *JsonUpdateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var mReq mdata.Metrics

	err := json.NewDecoder(req.Body).Decode(&mReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Wrong request body: %s", err), http.StatusBadRequest)
		return
	}

	err = j.saveData(&mReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Could not save data: %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

func (j *JsonUpdateHandler) saveData(mReq *mdata.Metrics) error {
	switch mReq.MType {
	case mdata.COUNTER:
		if mReq.Delta == nil {
			return fmt.Errorf("еmpty %s value", mReq.MType)
		}
		err := j.countStorage.Set(mdata.NewSimpleCounter(mReq.ID, *mReq.Delta))
		fmt.Println("Received: Counter", mReq.ID, *mReq.Delta)
		if err != nil {
			return fmt.Errorf("could not save data in storage")
		}
	case mdata.GAUGE:
		if mReq.Value == nil {
			return fmt.Errorf("еmpty %s value", mReq.MType)
		}
		err := j.gaugeStorage.Set(mdata.NewSimpleGauge(mReq.ID, *mReq.Value))
		fmt.Println("Received: Gauge", mReq.ID, *mReq.Value)
		if err != nil {
			return fmt.Errorf("could not save %s data in storage", mdata.GAUGE)
		}
	default:
		return fmt.Errorf("undefined metric type %s", mReq.MType)
	}

	return nil
}
