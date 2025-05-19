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

	routes := server.Routes{
		"/":                             shandler.NewGetListHandler(gaugeStorage, countStorage).ServeHTTP,
		"/update/{type}/{name}/{value}": shandler.NewUpdateHandler(mdata.InitMetrics(), gaugeStorage, countStorage).ServeHTTP,
		"/value/{type}/{name}":          shandler.NewGetHandler(mdata.InitMetrics(), gaugeStorage, countStorage).ServeHTTP,
		"/update/":                      shandler.NewJSONUpdateHandler(gaugeStorage, countStorage).ServeHTTP,
		"/value/":                       shandler.NewGetJSONMetricsHandler(gaugeStorage, countStorage).ServeHTTP,
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
		middlewares.NewCompressResponseMiddleware(),
		middlewares.NewLogRequestMiddleware(sugar),
		middlewares.NewDecompressRequestMiddleware(),
	}
}
