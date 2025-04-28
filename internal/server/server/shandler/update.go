package shandler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type UpdateRequest struct {
	Type  string
	Name  string
	Value string
}

func NewUpdateHandler(
	amt mdata.AvailableMetricsTypes,
	gaugeStorage server_storage.GaugeStorage,
	countStorage server_storage.CounterStorage,
) *UpdateHandler {
	return &UpdateHandler{
		AvailableMetricsTypes: amt,
		gaugeStorage:          gaugeStorage,
		countStorage:          countStorage,
	}
}

type UpdateHandler struct {
	AvailableMetricsTypes mdata.AvailableMetricsTypes
	gaugeStorage          server_storage.GaugeStorage
	countStorage          server_storage.CounterStorage
}

// TODO: как убрать дублирование writer.WriteHeader(http.StatusBadRequest)
func (uh *UpdateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	if req.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	err := uh.validateRequestHeader(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	ur, err := uh.updateRequestPrepare(req.URL.Path)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err = uh.validateUpdateRequest(ur); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.saveData(ur)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}

func (uh *UpdateHandler) saveData(ur *UpdateRequest) error {
	switch ur.Type {
	case mdata.COUNTER:
		val, err := strconv.ParseInt(ur.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s", mdata.COUNTER)
		}
		err = uh.countStorage.Set(mdata.NewSimpleCounter(ur.Name, val))
		fmt.Println("Received: Counter", ur.Name, val)
		if err != nil {
			return fmt.Errorf("could not save data in storage")
		}
	case mdata.GAUGE:
		val, err := strconv.ParseFloat(ur.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s", mdata.GAUGE)
		}
		err = uh.gaugeStorage.Set(mdata.NewSimpleGauge(ur.Name, val))
		fmt.Println("Received: Gauge", ur.Name, val)
		if err != nil {
			return fmt.Errorf("could not save %s data in storage", mdata.GAUGE)
		}
	default:
		return fmt.Errorf("undefined metric type %s", ur.Type)
	}

	return nil
}

// TODO: разделить на две функции - препейрер и создание?
func (uh *UpdateHandler) updateRequestPrepare(path string) (*UpdateRequest, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}
	ur := UpdateRequest{
		Type:  parts[2],
		Name:  parts[3],
		Value: parts[4],
	}
	return &ur, nil
}
