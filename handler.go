package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"private-ghp/config"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func setupHttpHandler() {
	config := config.GetConfig()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if page, ok := findPage(r, config); !ok {
			logrus.Debugf("page not found for request %s", r.Host)
			http.Error(w, "Page Not Found", http.StatusNotFound)
			return
		} else {
			logrus.Debugf("page found for request %s%s", r.Host, r.RequestURI)
			cookie, err := r.Cookie("token")
			if err != nil {
				redirectURL := fmt.Sprintf(
					"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=http://%s:%d/login/github/callback?origin=http://%s.%s:%d%s&scope=repo", config.Github.Client.Id, config.Domain, config.PublicPort, page.Subdomain, config.Domain, config.PublicPort, r.RequestURI)
				logrus.Debugf("no valid cookie found for request, redirecting to %s", redirectURL)
				http.Redirect(w, r, redirectURL, 301)
			} else {
				logrus.Debugf("cookie found for request %s", r.Host)
				ts := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: cookie.Value},
				)

				tc := oauth2.NewClient(oauth2.NoContext, ts)
				client := github.NewClient(tc)

				if r.RequestURI == "/" {
					r.RequestURI = "/" + page.Index
				}
				c, _, _, err := client.Repositories.GetContents(oauth2.NoContext, page.Repository.Owner, page.Repository.Name, r.RequestURI, &github.RepositoryContentGetOptions{Ref: page.Repository.Branch})
				logrus.Debugf("requesting content for owner: %s, repo: %s, path: %s, branch: %s", page.Repository.Owner, page.Repository.Name, r.RequestURI, page.Repository.Branch)
				if err == nil {
					logrus.Debugf("content found for owner: %s, repo: %s, path: %s, branch: %s", page.Repository.Owner, page.Repository.Name, r.RequestURI, page.Repository.Branch)
					sDec, _ := base64.StdEncoding.DecodeString(*c.Content)
					setContentType(c, w, r.RequestURI)

					setCache(page, c, w)
					w.Write(sDec)
				} else {
					logrus.Debugf("content not found for owner: %s, repo: %s, path: %s, branch: %s", page.Repository.Owner, page.Repository.Name, r.RequestURI, page.Repository.Branch)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

		}
	})

	http.HandleFunc("/login/github/callback", func(w http.ResponseWriter, r *http.Request) {
		origin := r.URL.Query().Get("origin")
		code := r.URL.Query().Get("code")
		token := getGithubAccessToken(code, config)

		w.Header().Add("Set-Cookie", fmt.Sprintf("token=%s; Domain=.%s; Path=/; HttpOnly", token, config.Domain))
		redirectURL := origin
		logrus.Debugf("token recevied from github (%s:%s), redirecting to %s", code, token, redirectURL)
		http.Redirect(w, r, redirectURL, 301)
	})
}

func findPage(r *http.Request, config *config.Config) (*config.Page, bool) {
	for _, page := range config.Pages {
		if fmt.Sprintf("%s.%s:%d", page.Subdomain, config.Domain, config.PublicPort) == r.Host || fmt.Sprintf("%s.%s", page.Subdomain, config.Domain) == r.Host {
			return &page, true
		}
	}
	return nil, false
}

func setCache(page *config.Page, c *github.RepositoryContent, w http.ResponseWriter) {
	w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", page.Cache.Duration))
	w.Header().Add("ETag", c.GetSHA())
	w.Header().Add("Expires", time.Now().Add(time.Duration(page.Cache.Duration)*time.Second).Format(http.TimeFormat))
	w.Header().Add("Last-Modified", time.Now().Format(http.TimeFormat))
}

func setContentType(c *github.RepositoryContent, w http.ResponseWriter, uri string) {
	t := "text/plain"

	if strings.HasSuffix(uri, ".html") || strings.HasSuffix(uri, ".htm") {
		t = "text/html"
	}
	if strings.HasSuffix(uri, ".css") {
		t = "text/css"
	}
	if strings.HasSuffix(uri, ".js") {
		t = "application/javascript"
	}
	if strings.HasSuffix(uri, ".jpg") || strings.HasSuffix(uri, ".jpeg") {
		t = "image/jpeg"
	}
	if strings.HasSuffix(uri, ".gif") {
		t = "image/gif"
	}
	if strings.HasSuffix(uri, ".png") {
		t = "image/png"
	}

	w.Header().Add("Content-Type", t)
}
