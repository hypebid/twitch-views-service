package rpc

import (
	"context"
	"runtime"

	"github.com/hypebid/go-kit/grpc/middleware"
	"github.com/hypebid/twitch-views-service/internal/config"
	"github.com/hypebid/twitch-views-service/internal/db"
	"github.com/hypebid/twitch-views-service/internal/rpc/pb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type Server struct {
	pb.UnsafeTwitchViewsServer
	Config *config.Config
}

var opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "twitchViews_processed_ops_total",
	Help: "The total number of processed HealthCheck events",
})

func initLogger(s *Server, tId string, methodName string) *logrus.Entry {
	// Build logger with TransactionId
	return s.Config.Log.WithFields(logrus.Fields{"transaction-id": tId, "method": methodName})
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthRequest) (*pb.HealthStatus, error) {
	// Build logger with TransactionId
	tId := ctx.Value(middleware.Grpc_ReqId_Marker)
	pc, _, _, _ := runtime.Caller(0)
	logger := initLogger(s, tId.(string), runtime.FuncForPC(pc).Name())

	logger.Info("received: ", req.GetMessage())

	// ping db
	dbOnline := false
	ping := db.PingDB(s.Config)
	if ping == nil {
		dbOnline = true
	}

	// add metric
	opsProcessed.Inc()

	return &pb.HealthStatus{
		TransactionId:  tId.(string),
		ServiceName:    s.Config.Constants.ServiceName,
		ReleaseDate:    s.Config.Constants.ReleaseDate,
		ReleaseSlug:    s.Config.Constants.ReleaseSlug,
		ReleaseVersion: s.Config.Constants.ReleaseVersion,
		DatabaseOnline: dbOnline,
		Message:        req.GetMessage(),
	}, nil
}

func (s *Server) GetStreamInfo(ctx context.Context, req *pb.TwitchUser) (*pb.StreamList, error) {
	// Build logger with TransactionId
	tId := ctx.Value(middleware.Grpc_ReqId_Marker)
	pc, _, _, _ := runtime.Caller(0)
	logger := initLogger(s, tId.(string), runtime.FuncForPC(pc).Name())

	logger.Info("hey there")

	return nil, nil
}
