module github.com/hypebid/twitch-views-service

go 1.16

require (
	github.com/go-chi/render v1.0.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hypebid/go-kit v0.1.3
	github.com/lib/pq v1.10.2
	github.com/prometheus/client_golang v1.11.0
	github.com/quasilyte/go-ruleguard/dsl v0.3.6
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	google.golang.org/grpc v1.39.1
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.12
)

// replace github.com/hypebid/go-kit => /Users/willmini/Development/go/src/github.com/hypebid/go-kit
