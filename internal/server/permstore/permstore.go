package permstore

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"ya-metrics/config"
	srvrstrg "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type PermStore struct {
	cfg      *config.PermStoreOptions
	storages []srvrstrg.Listable
	file     *os.File
	logger   *zap.SugaredLogger
}

type PermanentStorable interface {
	//метод достаёт данные из перманентного хранилища и помещает их в стораджи
	ExtractFromPermStore() error
	PutDataToPermStore() error
}

// TODO: прокинуть логгер
func NewPermStore(_ context.Context, logger *zap.SugaredLogger, cfg *config.PermStoreOptions, st ...srvrstrg.Listable) PermanentStorable {
	//открыть файл
	f, err := os.OpenFile(cfg.FileStoragePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	l := len(st)
	storages := make([]srvrstrg.Listable, l)
	for i, s := range st {
		storages[i] = s
	}
	return &PermStore{
		cfg:      cfg,
		storages: storages,
		file:     f,
		logger:   logger,
	}
}

func (ps *PermStore) ExtractFromPermStore() error {
	ps.logger.Info("ExtractFromPermStore")
	return nil
}

func (ps *PermStore) PutDataToPermStore() error {
	ps.logger.Info("PutDataToPermStore")
	var metrics []mdata.Metrics
	for _, s := range ps.storages {
		for _, t := range s.GetMetrics() {
			metrics = append(metrics, t)
		}
	}
	b, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	_, err = ps.file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
