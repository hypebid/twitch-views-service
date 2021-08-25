# go-micro-template
A Golang micro-service template. The boilerplate code for all of Hype Bid's micro-services coded in Go.

## Features
- gRPC Server
  - HealthCheck function
- gRPC Inteceptors / middleware
  - Tranaction Id
  - Logrus logger
  - ctx tags
  - Recover from panics
  - Hash validation check
- Gorm integration with Postgres
- [Doppler](https://docs.doppler.com/docs/start) for secrets
- Docker file
- Docker-Compose file
- Go modules
- Github Actions
- Prometheus metrics integrated
- Go tests
- [pre-commit ci](https://pre-commit.com/)

## Install pre-commit
```
brew install pre-commit
```

## Run Application Using Doppler
```
go build -o bin/micro-template -v .
doppler run --command="./bin/micro-template"
```

## Run Tests Using Doppler
```
doppler run --command="go test ./tests"
```