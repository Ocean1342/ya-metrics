package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"ya-metrics/config"
)

func NewChiServeable(cfg *config.Config, routes map[string]http.HandlerFunc) YaServeable {
	return &ChiServer{Config: cfg, routes: routes}
}

type ChiServer struct {
	Config *config.Config
	routes map[string]http.HandlerFunc
}

func (s *ChiServer) Start() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	for route, handler := range s.routes {
		router.HandleFunc(route, handler)
	}
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), router)
	if err != nil {
		panic(err)
	}
}
