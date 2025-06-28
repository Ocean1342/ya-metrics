package handlers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

func TestUpdateHandler_validateUpdateRequest(t *testing.T) {
	type args struct {
		ur *UpdateRequest
	}
	tests := []struct {
		name    string
		fields  Fields
		args    args
		wantErr bool
	}{
		{
			name:   "nil request",
			fields: initHandlerFields(),
			args: args{
				ur: nil,
			},
			wantErr: true,
		},
		{
			name:   "positive case",
			fields: initHandlerFields(),
			args: args{
				ur: &UpdateRequest{Type: mdata.GAUGE, Name: "Alloc", Value: "1.24"},
			},
			wantErr: false,
		},
		{
			name:   "negative case wrong name",
			fields: initHandlerFields(),
			args: args{
				ur: &UpdateRequest{Type: mdata.GAUGE, Name: "", Value: "1.24"},
			},
			wantErr: true,
		},
		{
			name:   "negative case wrong type",
			fields: initHandlerFields(),
			args: args{
				ur: &UpdateRequest{Type: "WrongType", Name: "some", Value: "1.24"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := initUpdateHandler(initHandlerFields())
			err := uh.validateUpdateRequest(tt.args.ur)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func initUpdateHandler(fields Fields) *Handler {
	return &Handler{
		availableMetricsTypes: fields.AvailableMetricsTypes,
		gaugeStorage:          fields.gaugeStorage,
		countStorage:          fields.countStorage,
	}
}

func initHandlerFields() Fields {
	f := Fields{
		AvailableMetricsTypes: mdata.InitMetrics(),
		gaugeStorage:          server_storage.NewSimpleGaugeStorage(),
		countStorage:          server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter),
	}
	return f
}

func TestUpdateHandler_validateMetricTypeName(t *testing.T) {
	tests := []struct {
		name    string
		fields  Fields
		ur      *UpdateRequest
		wantErr bool
	}{
		{
			name:    "positive case",
			fields:  initHandlerFields(),
			ur:      &UpdateRequest{Type: mdata.GAUGE, Name: "Alloc", Value: "1.24"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := initUpdateHandler(tt.fields)
			err := uh.validateMetricTypeName(tt.ur)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
