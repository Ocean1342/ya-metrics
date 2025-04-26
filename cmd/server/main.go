package main

import (
	"ya-metrics/config"
	"ya-metrics/internal/server/server"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/internal/server/server/shandler"
	"ya-metrics/pkg/mdata"
)

func main() {
	cfg := config.Config{
		Port: 8080,
		Host: "localhost",
	}

	gaugeStorage := server_storage.NewSimpleGaugeStorage()
	countStorage := server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter)

	updateHandler := shandler.NewUpdateHandler(mdata.InitMetrics(), gaugeStorage, countStorage)
	getHandler := shandler.NewGetHandler(mdata.InitMetrics(), gaugeStorage, countStorage)

	routes := server.Routes{
		"/update/{type}/{name}/{value}": updateHandler.HandlePost,
		"/value/{type}/{name}":          getHandler.ServeHTTP,
	}
	s := server.NewChiServeable(&cfg, routes)
	s.Start()
}
