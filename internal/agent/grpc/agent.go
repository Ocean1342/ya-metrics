package grpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"time"
	"ya-metrics/internal/agent/mgen"
	"ya-metrics/pkg/proto"
)

func Run(ctx context.Context, log *zap.SugaredLogger, reportIntervalSec int, target string) {
	l := log.With("grpc")
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ticker := time.NewTicker(time.Second * time.Duration(reportIntervalSec))
	defer ticker.Stop()
	c := proto.NewMetricsClient(conn)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			send(ctx, c, l)
		}
	}
}

func send(ctx context.Context, c proto.MetricsClient, log *zap.SugaredLogger) {
	for m := range mgen.GenerateGaugeMetrics(log) {
		ur := proto.UpdateMetricRequest{
			Metric: &proto.Metric{
				Name:  m.GetName(),
				Type:  m.GetType(),
				Value: fmt.Sprintf("%v", m.GetValue()),
			},
		}
		resp, err := c.UpdateMetric(ctx, &ur)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				log.Errorf("code:%s, err: %v, err text in response:%s", e.Code(), err, resp.Error)
			} else {
				log.Warnf("could not parse error, err: %v", err)
			}
		}
	}
}
