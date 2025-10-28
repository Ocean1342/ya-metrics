package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func (h *Handler) validateUpdateRequest(ur *UpdateRequest) error {
	if ur == nil {
		return fmt.Errorf("nil UpdateRequest")
	}
	err := h.validateMetricTypeName(ur)
	if err != nil {
		return err
	}
	//TODO: add validate name
	err = h.validateMetricName(ur)
	if err != nil {
		return err
	}
	//TODO: validate value?

	err = h.validateMetricValueNotEmpty(ur)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) validateMetricTypeName(ur *UpdateRequest) error {
	if !h.availableMetricsTypes.Isset(ur.Type) {
		return fmt.Errorf("no available metrics type found")
	}
	return nil
}

func (h *Handler) validateMetricName(ur *UpdateRequest) error {
	if strings.TrimSpace(ur.Name) == "" {
		return fmt.Errorf("no available metrics type found")
	}
	return nil
}

func (h *Handler) validateMetricValueNotEmpty(ur *UpdateRequest) error {
	if strings.TrimSpace(ur.Value) == "" {
		return fmt.Errorf("no available metrics type found")
	}
	return nil
}

func (h *Handler) validateRequestHeader(req *http.Request) error {
	//TODO: тест 3ий итерации
	//if strings.EqualFold(req.Header.Get("Content-Type"), "text/plain") {
	//	return nil
	//}
	//return fmt.Errorf("content type: %s not allowed", req.Header.Get("Content-Type"))
	return nil
}
