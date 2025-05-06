package main

import (
	"flag"
	"os"
	"ya-metrics/config"
	"ya-metrics/internal/server/server"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/internal/server/server/shandler"
	"ya-metrics/pkg/mdata"
)

func main() {
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
	s := server.NewChiServeable(&cfg, routes)
	s.Start()
}
