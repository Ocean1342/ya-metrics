package permstore

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"ya-metrics/config"
	"ya-metrics/internal/agent/mgen"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

func NewFullFilledPermStore() *PermStore {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("could not start logger")
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	permStoreOptions := config.PermStoreOptions{
		FileStoragePath: "./perm_storage.local.json",
		RestoreOnStart:  true,
		StoreInterval:   60,
	}
	gaugeStorage := server_storage.NewSimpleGaugeStorage()
	countStorage := server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)
	//todo: seed storages
	for _, m := range mgen.GenerateGaugeMetrics() {
		err = gaugeStorage.Set(m)
		if err != nil {
			errStr := fmt.Sprintf("could not put data to storage. Name: %s, type: %s, values:%f", m.GetName(), m.GetType(), m.GetValue())
			logger.Error(errStr)
		}
	}
	c := mdata.NewSimpleCounter("PollCount", 1)
	_ = countStorage.Set(c)
	return New(sugar, &permStoreOptions, gaugeStorage, countStorage)
}

func TestNewPermStore_PutDataToPermStore_PositiveCase(t *testing.T) {
	s := NewFullFilledPermStore()
	assert.NoError(t, s.Dump())
}

func TestPermStore_ExctractData(t *testing.T) {
	s := NewFullFilledPermStore()
	s.Dump()
	assert.NoError(t, s.Extract())
}

func TestPermStore_ExctractFrom(t *testing.T) {
	s := NewFullFilledPermStore()
	s.Dump()
	assert.NoError(t, s.Extract())
}
