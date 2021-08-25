package tests

import (
	"context"
	"os"
	"strings"

	"github.com/hypebid/twitch-views-service/internal/rpc/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *HeathCheckTestSuite) TestHealthCheckWithNoHash() {
	md := metadata.Pairs(
		"rpc-method", "healthCheck",
		"service-name", "testService",
		"hypebid-noauth", "false",
		"hypebid-nohash", "true",
		"hypebid-hash", "baldjfasdkfjkjsd",
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		s.T().Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewTwitchViewsClient(conn)
	resp, err := client.HealthCheck(ctx, &pb.HealthRequest{Message: "testing healthcheck"})
	if err != nil {
		s.T().Fatalf("health check failed: %v", err)
	}

	// asserts
	assert.Equal(s.T(), os.Getenv("SERVICE_NAME"), resp.ServiceName, "service name should match")
}

func (s *HeathCheckTestSuite) TestHealthCheckWithHash_Negative() {
	md := metadata.Pairs(
		"rpc-method", "healthCheck",
		"service-name", "testService",
		"hypebid-noauth", "false",
		"hypebid-nohash", "false",
		"hypebid-hash", "baldjfasdkfjkjsd",
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		s.T().Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewTwitchViewsClient(conn)
	resp, err := client.HealthCheck(ctx, &pb.HealthRequest{Message: "testing healthcheck"})

	// asserts
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), "auth issue")
}

func (s *HeathCheckTestSuite) TestHealthCheckWithHash_Positive() {
	md := metadata.Pairs(
		"rpc-method", "healthCheck",
		"service-name", "testService",
		"hypebid-noauth", "false",
		"hypebid-nohash", "false",
	)
	hash, err := createHash(md, strings.Split(s.Config.Constants.MetadataKeyList, ","), s.Config.Constants.HashSecret)
	if err != nil {
		s.T().Fatalf("failed to make hash: %v", err)
	}
	md.Append("hypebid-hash-bin", hash)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		s.T().Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewTwitchViewsClient(conn)
	resp, err := client.HealthCheck(ctx, &pb.HealthRequest{Message: "testing healthcheck"})
	if err != nil {
		s.T().Fatalf("failed to make request: %v", err)
	}

	// asserts
	assert.NotNil(s.T(), resp)
}
