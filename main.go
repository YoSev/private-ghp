package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"private-ghp/config"

	"github.com/sirupsen/logrus"
)

var configPath = flag.String("config", "/etc/private-ghp/config.yaml", "Path to config file yaml")

func main() {
	flag.Parse()

	err := config.Init(*configPath)
	if err != nil {
		logrus.Errorf("failed to load config: %s", err)
		os.Exit(1)
	}

	setupHttpHandler()

	logrus.Infof("listening on 0.0.0.0:%d", config.GetConfig().Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().Port), nil); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
