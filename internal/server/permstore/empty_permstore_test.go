package permstore

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"testing"
	"ya-metrics/config"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

const path = "./perm_storage.local.json"

func initEmptyPermStore(t *testing.T) PermanentStorable {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("could not start logger")
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Size() > 0 {
			err = os.Truncate(path, 0)
		}
		if err != nil {
			t.Fatalf("could not truncate file:%s; Error:%s", path, err)
		}
	}

	permStoreOptions := config.PermStoreOptions{
		FileStoragePath: path,
		RestoreOnStart:  true,
		StoreInterval:   60,
	}
	gaugeStorage := server_storage.NewSimpleGaugeStorage(sugar)
	countStorage := server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)
	return New(sugar, &permStoreOptions, gaugeStorage, countStorage)
}

func TestEmptyPutData(t *testing.T) {
	s := initEmptyPermStore(t)
	assert.NoError(t, s.Dump())
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	stat, err := f.Stat()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, int64(0), stat.Size())
}
