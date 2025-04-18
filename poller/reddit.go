package poller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type AuthResp struct {
	AccessToken string `json:"access_token"`
}

func (p *poller) setupReddit() (*http.Request, error) {
	data := url.Values{}
	userAgent := fmt.Sprintf("local:v0.0.1 (by /u/%s)", p.config.Username)
	data.Set("username", p.config.Username)
	data.Set("password", p.config.Password)
	data.Set("grant_type", "password")
	request, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("User-Agent", userAgent)
	request.SetBasicAuth(p.config.RedditAppID, p.config.RedditSecretKey)
	response, err := p.httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error while reading the response bytes:", err)
		return nil, err
	}
	var authResp AuthResp
	err = json.Unmarshal(respBody, &authResp)
	if err != nil {
		log.Fatal("Error while unmarshalling the response bytes:", err)
		return nil, err
	}

	bearer := "Bearer " + authResp.AccessToken
	newReq, err := http.NewRequest("GET", "https://oauth.reddit.com/r/funny/new/?t=hour&limit=50", nil)
	if err != nil {
		log.Fatal(err)
	}
	newReq.Header.Add("Authorization", bearer)
	newReq.Header.Add("User-Agent", userAgent)
	return newReq, nil
}
