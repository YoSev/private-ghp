package config

import (
	_ "embed"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port   int64  `yaml:"port"`
	Domain string `yaml:"domain"`
	Log    struct {
		Level string `yaml:"level"`
	}
	Github struct {
		Client struct {
			Id     string `yaml:"id"`
			Secret string `yaml:"secret"`
		}
	}
	Pages []Page
}

type Page struct {
	Subdomain string `yaml:"subdomain"`
	Index     string `yaml:"index"`
	Cache     struct {
		Duration int64 `yaml:"duration"`
	}
	Repository struct {
		Owner  string `yaml:"owner"`
		Name   string `yaml:"name"`
		Branch string `yaml:"branch"`
	}
}

var config *Config

//go:embed config.yaml
var defaultConfig []byte

func Init(configFile string) error {
	config = &Config{}

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	err := yaml.Unmarshal(defaultConfig, config)
	if err != nil {
		return err
	} else {
		setLogLevel()
	}

	if _, err := os.Stat(configFile); err == nil {
		cnt, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		} else {
			err = yaml.Unmarshal(cnt, config)
			if err != nil {
				return err
			} else {
				setLogLevel()
				logrus.Debugf("loaded config file: %s", configFile)
			}
		}
	} else {
		logrus.Debugf("config file '%s' not found, loaded defaults", configFile)
	}

	return nil
}

func GetConfig() *Config {
	return config
}

func setLogLevel() {
	switch strings.ToUpper(GetConfig().Log.Level) {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
		break
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
		break
	case "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
