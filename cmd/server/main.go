package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
)

var sugar *zap.SugaredLogger

func main() {
	initLogger()
	cfg := config.New()
	gaugeStorage := server_storage.NewSimpleGaugeStorage()
	countStorage := server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)
	//init perm store
	permStore := permstore.New(context.TODO(), sugar, cfg.PermStoreOptions, gaugeStorage, countStorage)
	if cfg.PermStoreOptions.RestoreOnStart {
		err := permStore.Extract()
		if err != nil {
			panic(fmt.Sprintf("panic on extract data from perm store on exit. err:%s", err))
		}
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//принудительная выгрузка при завершении работы
	go func() {
		v, ok := <-sigCh
		err := permStore.Put()
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
	}()

	h := handlers.New(gaugeStorage, countStorage, mdata.InitMetrics())
	routes := server.Routes{
		"/":                             h[handlers.GetListRoute].ServeHTTP,
		"/update/{type}/{name}/{value}": h[handlers.UpdateByURLParams].ServeHTTP,
		"/value/{type}/{name}":          h[handlers.GetByURLParams].ServeHTTP,
		"/update/":                      h[handlers.UpdateByJSON].ServeHTTP,
		"/value/":                       h[handlers.GetByJSON].ServeHTTP,
	}

	s := server.NewChiServeable(cfg, routes, initMiddlewares())
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

func initMiddlewares() []server.Middleware {
	return []server.Middleware{
		middlewares.NewLogResponseMiddleware(sugar),
		middlewares.NewCompressResponseMiddleware(),
		middlewares.NewLogRequestMiddleware(sugar),
		middlewares.NewDecompressRequestMiddleware(),
	}
}
