package rpc

import (
	"fmt"
	"net"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_reqAuth "github.com/hypebid/go-kit/grpc/middleware/auth"
	grpc_reqId "github.com/hypebid/go-kit/grpc/middleware/transactionId"
	"github.com/hypebid/twitch-views-service/internal/config"
	"github.com/hypebid/twitch-views-service/internal/rpc/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RpcSetup(c *config.Config) (net.Listener, *grpc.Server, error) {
	logOpts := []grpc_logrus.Option{}
	recovOpts := []grpc_recovery.Option{}
	reqAuthOpts := grpc_reqAuth.Options{
		HashSecret:      c.Constants.HashSecret,
		MetadataKeyList: strings.Split(c.Constants.MetadataKeyList, ","),
		MetadataHashKey: c.Constants.MetadataHashKey,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.Constants.Port))
	if err != nil {
		c.Log.Fatalf("Failed to listen: %v\n", err)
		return nil, nil, err
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithStreamServerChain(
			grpc_reqId.StreamServerInterceptor(c.Log),
			grpc_reqAuth.StreamServerInterceptor(c.Log),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(c.Log), logOpts...),
			grpc_recovery.StreamServerInterceptor(recovOpts...)),
		grpc_middleware.WithUnaryServerChain(
			grpc_reqId.UnaryServerInterceptor(c.Log),
			grpc_reqAuth.UnaryServerInterceptor(c.Log, reqAuthOpts),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(c.Log), logOpts...),
			grpc_recovery.UnaryServerInterceptor(recovOpts...)),
	)

	pb.RegisterServiceNameServer(grpcServer, &Server{Config: c})

	return lis, grpcServer, nil
}
