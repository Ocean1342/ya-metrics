package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
	cfg := config.New()
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
	go shutDown(permStore)
	privateCrypter, err := crypto.NewPrivateCrypter(cfg.CryptoKey, sugar)
	handler := handlers.New(gaugeStorage, countStorage, mdata.InitMetrics(), pg, sugar)
	s := server.NewChiServeable(cfg, handler, middlewares.InitMiddlewares(cfg, sugar, privateCrypter))
	s.Start()
}

func initLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("could not start logger")
	}
	defer logger.Sync()
	sugar = logger.Sugar()
}

func shutDown(permStore *permstore.PermStore) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//принудительная выгрузка при завершении работы
	v, ok := <-sigCh
	err := permStore.Dump()
	if err != nil {
		panic(fmt.Sprintf("panic on put data to perm store on exit. err:%s", err))
	}
	if ok {
		switch v {
		case syscall.SIGINT:
			os.Exit(int(syscall.SIGINT))
		case syscall.SIGTERM:
			os.Exit(int(syscall.SIGTERM))
		}
	}
}
