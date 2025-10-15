package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/pprof"
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
		"/updates/":                     s.handler.Updates,
		"/crypto":                       s.handler.CryptoGen,
	}
}

type ChiServer struct {
	Config      *config.Config
	handler     *handlers.Handler
	middlewares []Middleware
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
	err := http.ListenAndServe(s.Config.HostString, router)
	if err != nil {
		panic(err)
	}
}
