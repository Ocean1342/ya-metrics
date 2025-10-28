package handlers

import (
	"fmt"
	"net/http"
)

// GetList getting list of metrics
func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	var result string
	//TODO: вынести в отдельный интерфейс принтера
	for name, val := range h.gaugeStorage.GetList() {
		result += fmt.Sprintf("%s:%v<br>", name, val)
	}
	for name, val := range h.countStorage.GetList() {
		result += fmt.Sprintf("%s:%v", name, val)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
