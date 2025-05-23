package shandler

import (
	"fmt"
	"net/http"
	server_storage "ya-metrics/internal/server/server-storage"
)

type GetList struct {
	gaugeStorage server_storage.GaugeStorage
	countStorage server_storage.CounterStorage
}

func NewGetListHandler(gaugeStorage server_storage.GaugeStorage, countStorage server_storage.CounterStorage) *GetList {
	return &GetList{
		gaugeStorage: gaugeStorage,
		countStorage: countStorage,
	}
}

func (gl *GetList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	var result string
	//TODO: вынести в отдельный интерфейс принтера
	for name, val := range gl.gaugeStorage.GetList() {
		result += fmt.Sprintf("%s:%v<br>", name, val)
	}
	for name, val := range gl.countStorage.GetList() {
		result += fmt.Sprintf("%s:%v", name, val)
	}

	_, err := w.Write([]byte(result))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
