package server

import (
	"fmt"
	"net/http"
	"ya-metrics/config"
)

type YaServeable interface {
	Start()
}

type YaHttpServer struct {
	Config *config.Config
	routes map[string]HttpHandler
}

type Routes map[string]HttpHandler

type HttpHandler func(http.ResponseWriter, *http.Request)

func NewYaServeable(cfg *config.Config, routes map[string]HttpHandler) YaServeable {
	return &YaHttpServer{Config: cfg, routes: routes}
}

func (s *YaHttpServer) Start() {
	mux := http.NewServeMux()
	for route, handler := range s.routes {
		mux.HandleFunc(route, handler)
	}
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), mux)
	if err != nil {
		panic(err)
	}
}

//TODO: stop()
