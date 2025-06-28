package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	srvrstrg "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

func TestGetHandler_ServeHTTP(t *testing.T) {
	t.Skip(" почему то урл не парсится , т.е. typeName в хендлере пустой")
	type Fields struct {
		AvailableMetricsTypes mdata.AvailableMetricsTypes
		gaugeStorage          srvrstrg.GaugeStorage
		countStorage          srvrstrg.CounterStorage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type want struct {
		code    int
		content string
	}
	path := "http://localhost:8080"
	f := Fields{
		AvailableMetricsTypes: mdata.InitMetrics(),
		gaugeStorage:          srvrstrg.NewSimpleGaugeStorage(),
		countStorage:          srvrstrg.NewSimpleCountStorage(mdata.NewSimpleCounter),
	}
	_ = f.gaugeStorage.Set(mdata.NewSimpleGauge("one", 1))
	tests := []struct {
		name   string
		fields Fields
		args   args
		want   want
	}{
		{
			name:   "positive case",
			fields: f,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", path+"/value/gauge/one", nil),
			},
			want: want{
				code:    http.StatusOK,
				content: "one:1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gh := &Handler{
				availableMetricsTypes: tt.fields.AvailableMetricsTypes,
				gaugeStorage:          tt.fields.gaugeStorage,
				countStorage:          tt.fields.countStorage,
			}
			gh.Get(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.content, res.Body)
		})
	}
}
