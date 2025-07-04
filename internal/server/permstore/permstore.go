package permstore

import (
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
	Extract() error
	Dump() error
}

func New(logger *zap.SugaredLogger, cfg *config.PermStoreOptions, st ...srvrstrg.StorableStorage) *PermStore {
	//открыть файл
	f, err := os.OpenFile(cfg.FileStoragePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	l := len(st)
	storages := make([]srvrstrg.StorableStorage, l)
	copy(storages, st)
	permStore := PermStore{
		cfg:      cfg,
		storages: storages,
		file:     f,
		logger:   logger,
	}

	if cfg.StoreInterval != 0 {
		go func() {
			tick := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
			for range tick.C {
				if _, ok := <-tick.C; !ok {
					return
				}
				err = permStore.Dump()
				if err != nil {
					logger.Errorf("Error on put data to perm store on tick:%s", err)
				}
			}
		}()
	}
	if cfg.RestoreOnStart {
		err := permStore.Extract()
		if err != nil {
			panic(fmt.Sprintf("panic on extract data from perm store on exit. err:%s", err))
		}
	}
	return &permStore
}

func (ps *PermStore) Extract() error {
	ps.logger.Info("Extract")
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

func (ps *PermStore) Dump() error {
	ps.logger.Info("Dump")
	var metrics []mdata.Metrics
	err := ps.file.Truncate(0)
	if err != nil {
		return err
	}
	ps.file.Seek(0, io.SeekStart)
	for _, s := range ps.storages {
		metrics = append(metrics, s.GetMetrics()...)
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
