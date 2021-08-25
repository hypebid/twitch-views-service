package tests

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/hypebid/twitch-views-service/internal/config"
	"github.com/hypebid/twitch-views-service/internal/rpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type HeathCheckTestSuite struct {
	suite.Suite
	Config *config.Config
}

func init() {
	c, err := config.NewServiceConfig()
	if err != nil {
		log.Fatalf("failed to create config")
	}

	lis = bufconn.Listen(bufSize)
	_, grpcServer, err := rpc.RpcSetup(c)
	if err != nil {
		log.Fatalf("failed to setup rpc server")
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func (s *HeathCheckTestSuite) SetupTest() {
	c, err := config.NewServiceConfig()
	if err != nil {
		log.Fatal("failed to create config")
	}
	s.Config = c
}

func TestHealthCheckTestSuite(t *testing.T) {
	suite.Run(t, new(HeathCheckTestSuite))
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func createHash(md metadata.MD, mdKeyList []string, secret string) (string, error) {
	var hmac_message string
	for _, v := range mdKeyList {
		if len(md.Get(v)) == 0 {
			log.Printf("rpc request does not contain this metadata: %v", v)
			return "", errors.New("rpc request does not contain right metadata")
		}
		hmac_message = fmt.Sprintf("%v%v", hmac_message, md.Get(v)[0])
	}
	log.Println("hmac message created")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(hmac_message))
	expectedMAC := mac.Sum(nil)

	return string(expectedMAC), nil
}
