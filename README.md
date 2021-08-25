# twitch-views-service 

## Todo
- everything 

## Install pre-commit
```
brew install pre-commit
```

## Run Application Using Doppler
```
go build -o bin/twtich-views -v .
doppler run --command="./bin/twitch-views"
```

## Run Tests Using Doppler
```
doppler run --command="go test ./tests"
```