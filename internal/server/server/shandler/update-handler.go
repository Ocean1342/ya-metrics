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

func NewUpdateHandler(amt mdata.AvailableMetricsTypes) *UpdateHandler {
	return &UpdateHandler{AvailableMetricsTypes: amt}
}

type UpdateHandler struct {
	AvailableMetricsTypes mdata.AvailableMetricsTypes
	gaugeStorage          server_storage.GaugeStorage
	countStorage          server_storage.CounterStorage
}

func (uh *UpdateHandler) HandlePost(writer http.ResponseWriter, request *http.Request) {
	ur, err := uh.updateRequestPrepare(request.URL.Path)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err = uh.updateRequestValidateTypeName(ur); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: вынести в отдельную функцию сохранение метрики
	switch ur.Type {
	case mdata.COUNT:
		val, err := strconv.ParseInt(ur.Value, 10, 64)
		if err != nil {
			http.Error(writer, "invalid value for gauge", http.StatusBadRequest)
			return
		}
		err = uh.countStorage.Set(mdata.InitSimpleCounter(ur.Name, val))
		if err != nil {
			http.Error(writer, "invalid value for gauge", http.StatusBadRequest)
			return
		}
	case mdata.GAUGE:
		_, err := strconv.ParseFloat(ur.Value, 64)
		if err != nil {
			http.Error(writer, "invalid value for gauge", http.StatusBadRequest)
			return
		}
		uh.gaugeStorage.Get(ur.Name)
	}

}

func (uh *UpdateHandler) updateRequestValidateTypeName(ur *UpdateRequest) error {
	if !uh.AvailableMetricsTypes.Isset(ur.Type) {
		return fmt.Errorf("no available metrics type found")
	}
	//TODO: add validate name

	//TODO: validate value?

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
