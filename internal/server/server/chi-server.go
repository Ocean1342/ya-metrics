package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"ya-metrics/config"
	"ya-metrics/internal/server/server/handlers"
)

type Middleware func(http.Handler) http.Handler

func NewChiServeable(cfg *config.Config, handler *handlers.Handler, middlewares []Middleware, log *zap.SugaredLogger) YaServeable {
	return &ChiServer{Config: cfg, handler: handler, middlewares: middlewares, log: log}
}

func (s *ChiServer) initRoutes() map[string]http.HandlerFunc {
	return Routes{
		"/":                             s.handler.GetList,
		"/update/{type}/{name}/{value}": s.handler.Update,
		"/value/{type}/{name}":          s.handler.Get,
		"/update/":                      s.handler.UpdateByJSON,
		"/value/":                       s.handler.GetByJSON,
		"/ping":                         s.handler.Ping,
		"/updates/":                     s.handler.Updates,
		"/crypto":                       s.handler.CryptoGen,
	}
}

type ChiServer struct {
	Config      *config.Config
	log         *zap.SugaredLogger
	handler     *handlers.Handler
	middlewares []Middleware
	srv         *http.Server
}

func (s *ChiServer) Start() {
	router := chi.NewRouter()
	for _, m := range s.middlewares {
		router.Use(m)
	}
	if s.Config.ProfilingEnabled {
		router.Route("/debug/pprof", func(p chi.Router) {
			p.Get("/", pprof.Index)
			p.Get("/cmdline", pprof.Cmdline)
			p.Get("/profile", pprof.Profile)
			p.Post("/symbol", pprof.Symbol)
			p.Get("/symbol", pprof.Symbol)
			p.Get("/trace", pprof.Trace)
			p.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
			p.Get("/block", pprof.Handler("block").ServeHTTP)
			p.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
			p.Get("/heap", pprof.Handler("heap").ServeHTTP)
			p.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
			p.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		})
	}
	for route, handler := range s.initRoutes() {
		router.HandleFunc(route, handler)
	}
	srv := http.Server{
		Addr:    s.Config.HostString,
		Handler: router,
	}
	s.srv = &srv
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Errorf("server err: %v", err.Error())
	}
}

func (s *ChiServer) Stop(ctx context.Context) {
	s.log.Info("shutting down server")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Error("could not stop server")
	}
}
