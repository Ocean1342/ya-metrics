package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"ya-metrics/pkg/mdata"
)

type GetHandler struct {
	AvailableMetricsTypes mdata.AvailableMetricsTypes
}

// Get metric value
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	typeName := chi.URLParam(r, "type")
	if !h.availableMetricsTypes.Isset(typeName) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := chi.URLParam(r, "name")

	switch typeName {
	case mdata.GAUGE:
		g := h.gaugeStorage.Get(name)
		if g != nil {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(fmt.Sprintf("%v", g.GetValue())))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
	case mdata.COUNTER:
		c, err := h.countStorage.Get(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if c != nil {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(fmt.Sprintf("%v", c.GetValue())))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(""))
}
