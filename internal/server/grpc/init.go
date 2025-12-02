package grpc

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
	"ya-metrics/config"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
	"ya-metrics/pkg/proto"
)

func New(
	log *zap.SugaredLogger,
	cfg config.GRPCServerConfig,
	gaugeStorage server_storage.GaugeStorage,
	countStorage server_storage.CounterStorage,
	mTypes mdata.AvailableMetricsTypes,
) *MetricsServer {
	return &MetricsServer{
		log:                   log,
		cfg:                   cfg,
		gaugeStorage:          gaugeStorage,
		countStorage:          countStorage,
		availableMetricsTypes: mTypes,
	}
}

func (ms *MetricsServer) Run() {
	if !ms.cfg.Enabled {
		ms.log.Info("grpc server disabled")
		return
	}
	listen, err := net.Listen(ms.cfg.Network, ms.cfg.Addr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	ms.server = s
	proto.RegisterMetricsServer(s, ms)
	go func() {
		ms.log.Infof("start gRPC server on port: %s", ":3200")
		if err := s.Serve(listen); err != nil {
			ms.log.Fatalf("could not start gRPC server, err: %v", err)
		}
	}()
}

func (ms *MetricsServer) Stop() {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ms.server.GracefulStop()
	ms.log.Info("gRPC server stopped")
}
