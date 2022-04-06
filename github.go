package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"private-ghp/config"

	"github.com/sirupsen/logrus"
)

func getGithubAccessToken(code string, config *config.Config) string {
	requestBodyMap := map[string]string{
		"client_id":     config.Github.Client.Id,
		"client_secret": config.Github.Client.Secret,
		"code":          code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		logrus.Error("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		logrus.Error("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken
}
