package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"ya-metrics/config"
	"ya-metrics/internal/server/server/handlers"
)

type Middleware func(http.Handler) http.Handler

func NewChiServeable(cfg *config.Config, handler *handlers.Handler, middlewares []Middleware) YaServeable {
	return &ChiServer{Config: cfg, handler: handler, middlewares: middlewares}
}

func (s *ChiServer) initRoutes() map[string]http.HandlerFunc {
	return Routes{
		"/":                             s.handler.GetList,
		"/update/{type}/{name}/{value}": s.handler.Update,
		"/value/{type}/{name}":          s.handler.Get,
		"/update/":                      s.handler.UpdateByJSON,
		"/value/":                       s.handler.GetByJSON,
		"/ping":                         s.handler.Ping,
	}
}

type ChiServer struct {
	Config      *config.Config
	handler     *handlers.Handler
	middlewares []Middleware
}

func (s *ChiServer) Start() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	for _, m := range s.middlewares {
		router.Use(m)
	}
	for route, handler := range s.initRoutes() {
		router.HandleFunc(route, handler)
	}
	err := http.ListenAndServe(s.Config.HostString, router)
	if err != nil {
		panic(err)
	}
}
