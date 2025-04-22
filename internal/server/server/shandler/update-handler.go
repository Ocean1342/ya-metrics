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

func (uh *UpdateHandler) HandlePost(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	err := uh.validateRequestHeader(req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	ur, err := uh.updateRequestPrepare(req.URL.Path)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err = uh.validateUpdateRequest(ur); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.saveData(ur)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

}

func (uh *UpdateHandler) saveData(ur *UpdateRequest) error {
	switch ur.Type {
	case mdata.COUNT:
		val, err := strconv.ParseInt(ur.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for gauge")
		}
		err = uh.countStorage.Set(mdata.InitSimpleCounter(ur.Name, val))
		if err != nil {
			return fmt.Errorf("could not save data in storage")
		}
	case mdata.GAUGE:
		_, err := strconv.ParseFloat(ur.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for gauge")
		}
		uh.gaugeStorage.Get(ur.Name)
	default:
		return fmt.Errorf("undefined metric type")
	}

	return nil
}

// TODO: разделить на две функции - препейрер и создание?
func (uh *UpdateHandler) updateRequestPrepare(path string) (*UpdateRequest, error) {
	parts := strings.Split(path, "/")
	fmt.Println("len:", len(parts), cap(parts), parts)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}
	ur := UpdateRequest{}
	ur.Type = parts[2]
	ur.Name = parts[3]
	ur.Value = parts[4]
	return &ur, nil
}
