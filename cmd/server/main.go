package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"ya-metrics/config"
	"ya-metrics/internal/server/permstore"
	"ya-metrics/internal/server/server"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/internal/server/server-storage/database"
	"ya-metrics/internal/server/server/handlers"
	"ya-metrics/pkg/crypto"
	"ya-metrics/pkg/mdata"
	"ya-metrics/pkg/middlewares"
	"ya-metrics/pkg/postgres"
)

var (
	sugar        *zap.SugaredLogger
	buildVersion string = "N\\A"
	buildDate    string = "N\\A"
	buildCommit  string = "N\\A"
)

func main() {
	fmt.Printf("Build version:=%s, Build date=%s Build commit=%s\n", buildVersion, buildDate, buildCommit)
	initLogger()
	cfg := config.New(sugar)
	var permStore *permstore.PermStore
	var gaugeStorage server_storage.GaugeStorage
	var countStorage server_storage.CounterStorage
	pg, err := postgres.New(cfg.DBURL, sugar)
	if err != nil {
		sugar.Errorf("could not start pg. err: %s", err)
		gaugeStorage = server_storage.NewSimpleGaugeStorage(sugar)
		countStorage = server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)
		permStore = permstore.New(sugar, cfg.PermStoreOptions, gaugeStorage, countStorage)
	} else {
		gaugeStorage = database.NewGauge(pg, sugar)
		countStorage = database.NewCounter(pg, sugar, mdata.NewSimpleCounter)
		permStore = permstore.New(sugar, cfg.PermStoreOptions, gaugeStorage, countStorage)
	}
	privateCrypter, err := crypto.NewPrivateCrypter(cfg.CryptoKey, sugar)
	if err != nil {
		sugar.Errorf("could not create private crypter")
	}
	handler := handlers.New(gaugeStorage, countStorage, mdata.InitMetrics(), pg, sugar)
	s := server.NewChiServeable(cfg, handler, middlewares.InitMiddlewares(cfg, sugar, privateCrypter), sugar)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		shutDown(permStore, pg, s)
	}()
	s.Start()
	wg.Wait()
}

func shutDown(permStore *permstore.PermStore, pg *sql.DB, server server.YaServeable) {
	l := sugar.Named("graceful_shutdown")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	v, ok := <-sigCh
	l.Info("starting graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	server.Stop(ctx)
	cancel()
	err := permStore.Dump()
	if err != nil {
		panic(fmt.Sprintf("panic on put data to perm store on exit. err:%s", err))
	}
	if pg != nil {
		err = pg.Close()
		if err != nil {
			l.Warnf("could not close db conn. err: %v", err)
		}
	}
	if ok {
		switch v {
		case syscall.SIGINT:
			l.Info("finished by SIGINT")
		case syscall.SIGTERM:
			l.Info("finished by SIGTERM")
		case syscall.SIGQUIT:
			l.Info("finished by SIGQUIT")
		}
	}
}

func initLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("could not start logger")
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Printf("could not close logger")
		}
	}(logger)
	sugar = logger.Sugar()
}
