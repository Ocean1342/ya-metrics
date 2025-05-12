package main

import (
	"flag"
	"go.uber.org/zap"
	"os"
	"ya-metrics/config"
	"ya-metrics/internal/server/server"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/internal/server/server/shandler"
	"ya-metrics/pkg/mdata"
	"ya-metrics/pkg/middlewares"
)

var sugar *zap.SugaredLogger

func main() {
	initLogger()

	hostStr := flag.String("a", "localhost:8080", "server address")
	flag.Parse()
	if os.Getenv("ADDRESS") != "" {
		*hostStr = os.Getenv("ADDRESS")
	}

	cfg := config.Config{
		Port:       8080,
		Host:       "localhost",
		HostString: *hostStr,
	}

	gaugeStorage := server_storage.NewSimpleGaugeStorage()
	countStorage := server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)

	updateHandler := shandler.NewUpdateHandler(mdata.InitMetrics(), gaugeStorage, countStorage)
	getListHandler := shandler.NewGetListHandler(gaugeStorage, countStorage)
	getHandler := shandler.NewGetHandler(mdata.InitMetrics(), gaugeStorage, countStorage)

	routes := server.Routes{
		"/":                             getListHandler.ServeHTTP,
		"/update/{type}/{name}/{value}": updateHandler.ServeHTTP,
		"/value/{type}/{name}":          getHandler.ServeHTTP,
	}

	s := server.NewChiServeable(&cfg, routes, initMiddlewares())
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
		middlewares.NewLogRequestMiddleware(sugar),
	}
}
