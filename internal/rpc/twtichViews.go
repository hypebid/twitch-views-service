package rpc

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-chi/render"
	"github.com/hypebid/go-kit/grpc/middleware"
	"github.com/hypebid/twitch-views-service/internal/config"
	"github.com/hypebid/twitch-views-service/internal/db"
	"github.com/hypebid/twitch-views-service/internal/rpc/pb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type Server struct {
	pb.UnsafeTwitchViewsServer
	Config *config.Config
}

var token *oauth2.Token

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

func (s *Server) GetStreamInfo(ctx context.Context, req *pb.TwitchUser) (*pb.StreamInfo, error) {
	// Build logger with TransactionId
	tId := ctx.Value(middleware.Grpc_ReqId_Marker)
	pc, _, _, _ := runtime.Caller(0)
	logger := initLogger(s, tId.(string), runtime.FuncForPC(pc).Name())

	if token == nil {
		logger.Info("token is nil")
		err := getOAuthToken(logger, s.Config)
		if err != nil {
			return nil, err
		}
	}
	logger.Info("make request to get stream info")
	logger.Info("token: ", token.AccessToken)
	client := &http.Client{}
	r, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/streams?user_login=%v", req.GetUserLogin()), nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	r.Header.Add("Client-id", s.Config.Constants.TwitchClientId)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	logger.Info("request done successfully")
	var respMap StreamInfo
	err = render.DecodeJSON(resp.Body, &respMap)
	if err != nil {
		return nil, err
	}
	logger.Info("resp mapped")

	return &pb.StreamInfo{
		TransactionId: tId.(string),
		StreamId:      respMap.Data[0].ID,
		UserId:        respMap.Data[0].UserID,
		UserLogin:     respMap.Data[0].UserLogin,
		UserName:      respMap.Data[0].UserName,
		IsLive:        true,
		ViewerCount:   int32(respMap.Data[0].ViewerCount),
		StartedAt:     respMap.Data[0].StartedAt.String(),
		Language:      respMap.Data[0].Language,
		IsMature:      respMap.Data[0].IsMature,
	}, nil
}

func getOAuthToken(logger *logrus.Entry, config *config.Config) error {
	logger.Info("getting twtich OAuth token")

	oauth2Config := &clientcredentials.Config{
		ClientID:     config.Constants.TwitchClientId,
		ClientSecret: config.Constants.TwitchSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	t, err := oauth2Config.Token(context.Background())
	if err != nil {
		logger.Error("error getting twitch OAuth token")
		return err
	}
	token = t

	return nil
}
