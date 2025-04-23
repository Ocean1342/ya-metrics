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

	updateHandler := shandler.NewUpdateHandler(
		mdata.InitMetrics(),
		server_storage.NewSimpleGaugeStorage(),
		server_storage.NewSimpleCountStorage(mdata.NewSimpleCounter),
	)
	routes := server.Routes{
		"/update/{type}/{name}/{value}": updateHandler.HandlePost,
	}
	s := server.NewYaServeable(&cfg, routes)
	s.Start()
}
