package permstore

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"time"
	"ya-metrics/config"
	srvrstrg "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
)

type PermStore struct {
	cfg      *config.PermStoreOptions
	storages []srvrstrg.StorableStorage
	file     *os.File
	logger   *zap.SugaredLogger
}

type PermanentStorable interface {
	//метод достаёт данные из перманентного хранилища и помещает их в стораджи
	ExtractFromPermStore() error
	PutDataToPermStore() error
}

func NewPermStore(_ context.Context, logger *zap.SugaredLogger, cfg *config.PermStoreOptions, st ...srvrstrg.StorableStorage) PermanentStorable {
	//открыть файл
	f, err := os.OpenFile(cfg.FileStoragePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	l := len(st)
	storages := make([]srvrstrg.StorableStorage, l)
	for i, s := range st {
		storages[i] = s
	}
	permStore := PermStore{
		cfg:      cfg,
		storages: storages,
		file:     f,
		logger:   logger,
	}

	if cfg.StoreInterval != 0 {
		go func() {
			tick := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
			for {
				select {
				case <-tick.C:
					err = permStore.PutDataToPermStore()
					if err != nil {
						logger.Errorf("Error on put data to perm store on tick:%s", err)
					}
				}
			}
		}()
	}

	return &permStore
}

func (ps *PermStore) ExtractFromPermStore() error {
	ps.logger.Info("ExtractFromPermStore")
	var metrics []mdata.Metrics
	//смаршалить метрики
	bytes, err := io.ReadAll(ps.file)
	if err != nil {
		return err
	}
	if len(bytes) == 0 {
		return nil
	}

	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		return err
	}
	//запихать их в сторадж
	for _, s := range ps.storages {
		err = s.SetFrom(metrics)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PermStore) PutDataToPermStore() error {
	ps.logger.Info("PutDataToPermStore")
	var metrics []mdata.Metrics
	err := ps.file.Truncate(0)
	if err != nil {
		return err
	}
	ps.file.Seek(0, io.SeekStart)
	for _, s := range ps.storages {
		for _, t := range s.GetMetrics() {
			metrics = append(metrics, t)
		}
	}
	size := len(metrics)
	ps.logger.Info(fmt.Sprintf("Metrics len:%d", size))
	if size != 0 {
		b, err := json.Marshal(metrics)
		if err != nil {
			return err
		}
		_, err = ps.file.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}
