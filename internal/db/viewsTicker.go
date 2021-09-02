package db

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/render"
	"github.com/hypebid/twitch-views-service/internal/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

var token *oauth2.Token
var retries int

func TwitchViewsTicker(c *config.Config, userString string) {
	pc, _, _, _ := runtime.Caller(0)
	logger := c.Log.WithFields(logrus.Fields{"method": runtime.FuncForPC(pc).Name()})

	if token == nil {
		logger.Info("token is nil")
		err := getOAuthToken(logger, c)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	for {
		logger.Info("make request to get stream info")
		logger.Debug("token: ", token.AccessToken)

		respMap, err := makeRequest(logger, c, userString)
		if err != nil {
			logger.Error(err)
			return
		}

		logger.Info("======================================")
		logger.Printf("request complete: %+v", respMap)
		logger.Info("======================================")

		time.Sleep(time.Second * 15)
	}

}

func makeRequest(logger *logrus.Entry, c *config.Config, stream string) (*StreamInfo, error) {
	client := &http.Client{}
	r, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/streams?user_login=%v", stream), nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	r.Header.Add("Client-id", c.Constants.TwitchClientId)
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

	// check that a value was returned
	if len(respMap.DataList) < 1 {
		logger.Info("entering retry logic")
		// retry logic
		if retries < 3 {
			retries++
			err := getOAuthToken(logger, c)
			if err != nil {
				return nil, err
			}
			logger.Info("doing retry # ", retries)
			// recurrsion
			//nolint
			return makeRequest(logger, c, stream)
		}
		retries = 0
		logger.Error(errors.New("ran out of retries. cancel request"))
		return nil, errors.New("stream offline / doesn't exist")
	}

	return &respMap, nil
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
