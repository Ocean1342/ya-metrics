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
	"ya-metrics/internal/server/server/shandler"
	"ya-metrics/pkg/mdata"
	"ya-metrics/pkg/middlewares"
)

var sugar *zap.SugaredLogger

//TODO:
//	    -1 - реализовать выгрузку по завершению приложения, done
//		0 - продумать передачу метрик при синхронном режиме
//		1 - реализация новых хранилищ
//		2 - реализация интерфейса PermanentStorable

func main() {
	initLogger()
	cfg := config.InitConfig()
	gaugeStorage := server_storage.NewSimpleGaugeStorage()
	countStorage := server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)
	//init perm store
	permStore := permstore.NewPermStore(context.TODO(), sugar, cfg.PermStoreOptions, gaugeStorage, countStorage)
	if cfg.PermStoreOptions.RestoreOnStart {
		err := permStore.ExtractFromPermStore()
		if err != nil {
			panic(fmt.Sprintf("panic on extract data from perm store on exit. err:%s", err))
		}
	}
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//принудительная выгрузка при завершении работы
	go func() {
		for {
			select {
			case v, ok := <-sigCh:
				err := permStore.PutDataToPermStore()
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
		}
	}()

	routes := server.Routes{
		"/":                             shandler.NewGetListHandler(gaugeStorage, countStorage).ServeHTTP,
		"/update/{type}/{name}/{value}": shandler.NewUpdateHandler(mdata.InitMetrics(), gaugeStorage, countStorage).ServeHTTP,
		"/value/{type}/{name}":          shandler.NewGetHandler(mdata.InitMetrics(), gaugeStorage, countStorage).ServeHTTP,
		"/update/":                      shandler.NewJSONUpdateHandler(gaugeStorage, countStorage).ServeHTTP,
		"/value/":                       shandler.NewGetJSONMetricsHandler(gaugeStorage, countStorage).ServeHTTP,
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
