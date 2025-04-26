package shandler

import (
	"fmt"
	"net/http"
	"strings"
)

func (uh *UpdateHandler) validateUpdateRequest(ur *UpdateRequest) error {
	if ur == nil {
		return fmt.Errorf("nil UpdateRequest")
	}
	err := uh.validateMetricTypeName(ur)
	if err != nil {
		return err
	}
	//TODO: add validate name
	err = uh.validateMetricName(ur)
	if err != nil {
		return err
	}
	//TODO: validate value?

	err = uh.validateMetricValueNotEmpty(ur)
	if err != nil {
		return err
	}
	return nil
}

func (uh *UpdateHandler) validateMetricTypeName(ur *UpdateRequest) error {
	if !uh.AvailableMetricsTypes.Isset(ur.Type) {
		return fmt.Errorf("no available metrics type found")
	}
	return nil
}

func (uh *UpdateHandler) validateMetricName(ur *UpdateRequest) error {
	if strings.TrimSpace(ur.Name) == "" {
		return fmt.Errorf("no available metrics type found")
	}
	return nil
}

func (uh *UpdateHandler) validateMetricValueNotEmpty(ur *UpdateRequest) error {
	if strings.TrimSpace(ur.Value) == "" {
		return fmt.Errorf("no available metrics type found")
	}
	return nil
}

func (uh *UpdateHandler) validateRequestHeader(req *http.Request) error {
	if strings.EqualFold(req.Header.Get("Content-Type"), "text/plain") {
		return nil
	}
	return fmt.Errorf("content type: %s not allowed", req.Header.Get("Content-Type"))
}
