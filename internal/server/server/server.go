package server

import (
	"context"
	"net/http"
)

type YaServeable interface {
	Start()
	Stop(ctx context.Context)
}

type Routes map[string]http.HandlerFunc
