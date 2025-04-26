package server

import (
	"fmt"
	"net/http"
	"ya-metrics/config"
)

type YaServeable interface {
	Start()
}

type YaHTTPServer struct {
	Config *config.Config
	routes map[string]http.HandlerFunc
}

type Routes map[string]http.HandlerFunc

func NewYaServeable(cfg *config.Config, routes map[string]http.HandlerFunc) YaServeable {
	return &YaHTTPServer{Config: cfg, routes: routes}
}

func (s *YaHTTPServer) Start() {
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
