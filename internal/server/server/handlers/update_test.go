package handlers

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type want struct {
	code        int
	contentType string
}
type Fields struct {
	AvailableMetricsTypes mdata.AvailableMetricsTypes
	gaugeStorage          server_storage.GaugeStorage
	countStorage          server_storage.CounterStorage
}
type reqParams struct {
	method string
	path   string
}

// TODO: Вопрос: нужно ли в тесте хендлера проверять запись в сторадж?
func TestUpdateHandler(t *testing.T) {
	path := "http://localhost:8080"
	logger, _ := zap.NewProduction()
	f := Fields{
		AvailableMetricsTypes: mdata.InitMetrics(),
		gaugeStorage:          server_storage.NewSimpleGaugeStorage(logger.Sugar()),
		countStorage:          server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter),
	}
	tests := []struct {
		name      string
		fields    Fields
		reqParams reqParams
		want      want
	}{
		{
			name:   "positive case gauge",
			fields: f,
			reqParams: reqParams{
				method: http.MethodPost,
				path:   path + "/update/gauge/Alloc/1.23456",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:   "positive case count",
			fields: f,
			reqParams: reqParams{
				method: http.MethodPost,
				path:   path + "/update/counter/PollCount/1",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:   "negative case wrong type",
			fields: f,
			reqParams: reqParams{
				method: http.MethodPost,
				path:   path + "/update/wrong-type/PollCount/1",
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//создать реквест
			req := httptest.NewRequest(tt.reqParams.method, tt.reqParams.path, nil)
			req.Header.Set("Content-Type", "text/plain")
			//новый ридер
			nr := httptest.NewRecorder()
			//дёрнуть хендлер
			uh := &Handler{
				availableMetricsTypes: tt.fields.AvailableMetricsTypes,
				gaugeStorage:          tt.fields.gaugeStorage,
				countStorage:          tt.fields.countStorage,
			}
			uh.Update(nr, req)
			//nr получить результ
			res := nr.Result()
			defer res.Body.Close()
			//сравнить результаты
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
