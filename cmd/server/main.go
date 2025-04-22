package main

import (
	"ya-metrics/config"
	"ya-metrics/internal/server/server"
	"ya-metrics/internal/server/server/shandler"
	"ya-metrics/pkg/mdata"
)

func main() {
	//TODO: вынести конфиг на энвы?
	cfg := config.Config{
		Port: 8080,
		Host: "localhost",
	}
	routes := server.Routes{
		"/update/{type}/{name}/{value}": (shandler.NewUpdateHandler(mdata.InitMetrics())).HandlePost,
	}
	s := server.NewYaServeable(&cfg, routes)
	s.Start()
}
