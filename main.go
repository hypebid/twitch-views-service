package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/quasilyte/go-ruleguard/dsl"

	"github.com/hypebid/twitch-views-service/internal/config"
	"github.com/hypebid/twitch-views-service/internal/rpc"
)

func metrics(c *config.Config) {
	http.Handle("/metrics", promhttp.Handler())
	c.Log.Fatal(http.ListenAndServe(":2111", nil))
}

func main() {
	c, err := config.NewServiceConfig()
	if err != nil {
		log.Printf("Error initializing service config: %v", err)
		return
	}

	c.Log.Info("starting metrics route...")
	go metrics(c)
	c.Log.Info("done")

	c.Log.Info("setting up grpc server...")
	lis, grpcServer, err := rpc.RpcSetup(c)
	if err != nil {
		c.Log.Error("error setting up rpc server: ", err)
		return
	}
	c.Log.Info("done")

	c.Log.Info("Server listening on ", c.Constants.Port)

	if err := grpcServer.Serve(lis); err != nil {
		c.Log.Fatalf("failed to serve: %v\n", err)
	}
}
