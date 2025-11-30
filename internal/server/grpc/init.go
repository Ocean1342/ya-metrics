package grpc

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ya-metrics/config"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/pkg/mdata"
	"ya-metrics/pkg/proto"
)

func New(
	gaugeStorage server_storage.GaugeStorage,
	countStorage server_storage.CounterStorage,
	mTypes mdata.AvailableMetricsTypes,
) *MetricsServer {
	return &MetricsServer{
		gaugeStorage:          gaugeStorage,
		countStorage:          countStorage,
		availableMetricsTypes: mTypes,
	}
}

func Init(
	sugar *zap.SugaredLogger,
	cfg config.GPRCServerConfig,
	gaugeStorage server_storage.GaugeStorage,
	countStorage server_storage.CounterStorage,
	mTypes mdata.AvailableMetricsTypes) {
	if !cfg.Enabled {
		sugar.Info("grpc server disabled")
		return
	}
	server := New(gaugeStorage, countStorage, mTypes)
	listen, err := net.Listen(cfg.Network, cfg.Addr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	proto.RegisterMetricsServer(s, server)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sugar.Infof("start gRPC server on port: %s", ":3200")
		if err := s.Serve(listen); err != nil {
			sugar.Fatalf("could not start gRPC server, err: %v", err)
		}
	}()
	<-stop
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.GracefulStop()
	sugar.Info("gRPC server stopped")
}
