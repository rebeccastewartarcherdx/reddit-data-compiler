package poller

import (
	"encoding/json"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"io"
	"log"
	"net/http"
	"net/url"
	"redditDataCompiler/models"
	"strconv"
	"time"
)

type Poller interface {
	Poll() chan PollResponse
}

type poller struct {
	httpClient *http.Client
	config     envConfig
}

type PollResponse struct {
	Data *models.APIResponse
	Err  error
}

func NewPoller(httpClient *http.Client) Poller {
	var config envConfig
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	return &poller{
		httpClient: httpClient,
		config:     config,
	}
}

type envConfig struct {
	Username        string `envconfig:"REDDIT_USERNAME"`
	Password        string `envconfig:"REDDIT_PASSWORD"`
	RedditAppID     string `envconfig:"REDDIT_APPID"`
	RedditSecretKey string `envconfig:"REDDIT_SECRET"`
}

var (
	apiLimitExceeded = errors.New("API limit exceeded")
)

func (p *poller) Poll() chan PollResponse {
	request, err := p.setupReddit()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	respChan := make(chan PollResponse)
	// Goroutine that periodically pulls data from the API.
	go func() {
		defer func() {
			close(respChan)
		}()
		ticker := time.Tick(2 * time.Second) // fetch data every 2 seconds
		for range ticker {
			apiResponse, err := p.fetchRedditData(request, "https://oauth.reddit.com/r/funny/new/?t=hour&limit=50")
			if err != nil {
				respChan <- PollResponse{
					Data: nil,
					Err:  err,
				}
				break
			}
			respChan <- PollResponse{apiResponse, nil}
		}
	}()
	return respChan
}

// putting url here to make it easier to unit test
func (p *poller) fetchRedditData(request *http.Request, urlString string) (*models.APIResponse, error) {
	var err error
	request.URL, err = url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	resp, err := p.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	// too many requests
	if resp.StatusCode == 429 || resp.Header.Get("X-RateLimit-Remaining") == "0" {
		sleepTime, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Reset"))
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)
		return nil, apiLimitExceeded
	}

	// parse response and send to channel to process
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var apiResponse models.APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
