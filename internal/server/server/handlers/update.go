package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"ya-metrics/pkg/mdata"
)

type UpdateRequest struct {
	Type  string
	Name  string
	Value string
}

func (h *Handler) Update(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	if req.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	err := h.validateRequestHeader(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	ur, err := UpdateRequestPrepare(req.URL.Path)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.validateUpdateRequest(ur); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.saveData(ur)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}

func (h *Handler) saveData(ur *UpdateRequest) error {
	switch ur.Type {
	case mdata.COUNTER:
		val, err := strconv.ParseInt(ur.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s", mdata.COUNTER)
		}
		err = h.countStorage.Set(mdata.NewSimpleCounter(ur.Name, val))
		fmt.Println("Received: Counter", ur.Name, val)
		if err != nil {
			return fmt.Errorf("could not save data in storage")
		}
	case mdata.GAUGE:
		val, err := strconv.ParseFloat(ur.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s", mdata.GAUGE)
		}
		err = h.gaugeStorage.Set(mdata.NewSimpleGauge(ur.Name, val))
		fmt.Println("Received: Gauge", ur.Name, val)
		if err != nil {
			return fmt.Errorf("could not save %s data in storage", mdata.GAUGE)
		}
	default:
		return fmt.Errorf("undefined metric type %s", ur.Type)
	}

	return nil
}

// UpdateRequestPrepare - prepare request from url path
func UpdateRequestPrepare(path string) (*UpdateRequest, error) {
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
