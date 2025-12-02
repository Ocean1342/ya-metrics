package grpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"ya-metrics/config"
	server_storage "ya-metrics/internal/server/server-storage"
	"ya-metrics/internal/server/server/handlers"
	"ya-metrics/pkg/mdata"
	"ya-metrics/pkg/proto"
)

type MetricsServer struct {
	log *zap.SugaredLogger
	cfg config.GRPCServerConfig
	// availableMetricsTypes - available metrics types e.g. gauge, counter
	availableMetricsTypes mdata.AvailableMetricsTypes
	// gaugeStorage - gauge type storage
	gaugeStorage server_storage.GaugeStorage
	//countStorage - count type storage
	countStorage server_storage.CounterStorage
	proto.UnimplementedMetricsServer
	server *grpc.Server
}

func (ms *MetricsServer) UpdateMetric(_ context.Context, req *proto.UpdateMetricRequest) (*proto.UpdateMetricResponse, error) {
	resp := &proto.UpdateMetricResponse{}
	ur := &handlers.UpdateRequest{
		Type:  req.Metric.Type,
		Name:  req.Metric.Name,
		Value: req.Metric.Value,
	}
	err := ms.saveData(ur)
	if err != nil {
		resp.Error = fmt.Sprintf("could not save data: %s", err)
	} else {
		ms.log.Infof("metric updated: %s by grpc", ur.Name)
	}
	return resp, nil
}

func (ms *MetricsServer) GetMetric(_ context.Context, req *proto.GetMetricRequest) (*proto.GetMetricResponse, error) {
	resp := &proto.GetMetricResponse{}
	if !ms.availableMetricsTypes.Isset(req.GetType()) {
		resp.Error = fmt.Sprintf("%s is not available metric type", req.GetType())
		return resp, nil
	}

	switch req.GetType() {
	case mdata.GAUGE:
		g := ms.gaugeStorage.Get(req.GetName())
		if g != nil {
			resp.Metric = &proto.Metric{
				Type:  g.GetType(),
				Name:  g.GetName(),
				Value: strconv.FormatFloat(g.GetValue(), 'g', -1, 64),
			}
			return resp, nil
		}
		return nil, status.Errorf(codes.NotFound, "metric type: %s name: %s not found", req.GetType(), req.GetName())
	case mdata.COUNTER:
		c, err := ms.countStorage.Get(req.GetName())
		if err != nil {
			resp.Error = fmt.Sprintf("get counter err: %v", err)
			return resp, nil
		}
		if c != nil {
			resp.Metric = &proto.Metric{
				Type:  c.GetType(),
				Name:  c.GetName(),
				Value: strconv.Itoa(int(c.GetValue())),
			}
			return resp, nil
		}
		return nil, status.Errorf(codes.NotFound, "metric type: %s name: %s not found", req.GetType(), req.GetName())
	}
	return nil, status.Errorf(codes.Internal, "metric type: %s name: %s not found", req.GetType(), req.GetName())
}
func (ms *MetricsServer) saveData(ur *handlers.UpdateRequest) error {
	switch ur.Type {
	case mdata.COUNTER:
		val, err := strconv.ParseInt(ur.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s", mdata.COUNTER)
		}
		err = ms.countStorage.Set(mdata.NewSimpleCounter(ur.Name, val))
		fmt.Println("Received: Counter", ur.Name, val)
		if err != nil {
			return fmt.Errorf("could not save data in storage")
		}
	case mdata.GAUGE:
		val, err := strconv.ParseFloat(ur.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s", mdata.GAUGE)
		}
		err = ms.gaugeStorage.Set(mdata.NewSimpleGauge(ur.Name, val))
		fmt.Println("Received: Gauge", ur.Name, val)
		if err != nil {
			return fmt.Errorf("could not save %s data in storage", mdata.GAUGE)
		}
	default:
		return fmt.Errorf("undefined metric type %s", ur.Type)
	}

	return nil
}
