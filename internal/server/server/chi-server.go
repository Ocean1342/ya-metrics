package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"ya-metrics/config"
)

type Middleware func(http.Handler) http.Handler

func NewChiServeable(cfg *config.Config, routes map[string]http.HandlerFunc, middlewares []Middleware) YaServeable {
	return &ChiServer{Config: cfg, routes: routes, middlewares: middlewares}
}

type ChiServer struct {
	Config      *config.Config
	routes      map[string]http.HandlerFunc
	middlewares []Middleware
}

func (s *ChiServer) Start() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	for _, m := range s.middlewares {
		router.Use(m)
	}
	for route, handler := range s.routes {
		router.HandleFunc(route, handler)
	}
	err := http.ListenAndServe(s.Config.HostString, router)
	if err != nil {
		panic(err)
	}
}
