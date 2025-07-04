package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ya-metrics/config"
	"ya-metrics/internal/server/permstore"
	"ya-metrics/internal/server/server"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/internal/server/server/handlers"
	"ya-metrics/pkg/mdata"
	"ya-metrics/pkg/middlewares"
	"ya-metrics/pkg/postgres"
)

var sugar *zap.SugaredLogger

func main() {
	initLogger()
	cfg := config.New()
	pg, err := postgres.New(cfg.DBURL)
	if err != nil {
		log.Fatalf("could not start db.err:%s", err)
	}
	gaugeStorage := server_storage.NewSimpleGaugeStorage()
	countStorage := server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)
	//init perm store
	permStore := permstore.New(sugar, cfg.PermStoreOptions, gaugeStorage, countStorage)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//принудительная выгрузка при завершении работы
	go func() {
		v, ok := <-sigCh
		err := permStore.Dump()
		if err != nil {
			panic(fmt.Sprintf("panic on put data to perm store on exit. err:%s", err))
		}
		if ok {
			switch v {
			case syscall.SIGINT:
				pg.Close()
				os.Exit(int(syscall.SIGINT))
			case syscall.SIGTERM:
				pg.Close()
				os.Exit(int(syscall.SIGTERM))
			}
		}
	}()

	handler := handlers.New(gaugeStorage, countStorage, mdata.InitMetrics(), pg)
	s := server.NewChiServeable(cfg, handler, middlewares.InitMiddlewares(sugar))
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
