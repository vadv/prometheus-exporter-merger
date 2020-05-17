package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	defaultListen        = ":8080"
	defaultScrapeTimeout = 15 * time.Second
)

type source struct {
	Url    string            `yaml:"url"`
	Labels map[string]string `yaml:"labels"`
}

type config struct {
	Listen        string        `yaml:"listen"`
	ScrapeTimeout time.Duration `yaml:"scrape_timeout"`
	Sources       []*source     `yaml:"sources"`
}

func parseConfig(filename string) (*config, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return parseConfigFromFile(filename)
	}
	if os.IsNotExist(err) {
		return parseConfigFromEnv()
	}
	return nil, err
}

func parseConfigFromFile(filename string) (*config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result := &config{
		Listen:        defaultListen,
		ScrapeTimeout: defaultScrapeTimeout,
	}
	return result, yaml.Unmarshal(data, result)
}

func parseConfigFromEnv() (*config, error) {
	result := &config{
		Listen:        defaultListen,
		ScrapeTimeout: defaultScrapeTimeout,
	}
	if v := os.Getenv("LISTEN"); v != "" {
		result.Listen = v
	}
	if v := os.Getenv("SCRAPE_TIMEOUT"); v != "" {
		timeout, err := time.ParseDuration(v)
		if err != nil {
			return nil, errors.Wrap(err, "parse SCRAPE_TIMEOUT")
		}
		result.ScrapeTimeout = timeout
	}
	for _, env := range os.Environ() {
		// URL_ONE=http://127.0.0.1:8080/metrics,k1:v1,k2:v2
		if strings.HasPrefix(env, "URL_") {
			args := strings.Split(env, "=")
			if len(args) != 2 {
				return nil, fmt.Errorf("unable to parse env variable %s", env)
			}
			valuesArgs := strings.Split(args[1], ",")
			s := &source{Url: valuesArgs[0], Labels: make(map[string]string)}
			if len(valuesArgs) > 1 {
				for i, v := range valuesArgs {
					if i == 0 {
						continue
					}
					labelArgs := strings.Split(v, ":")
					if len(labelArgs) != 2 {
						return nil, fmt.Errorf("unable to parse labels from env variable %s", env)
					}
					s.Labels[labelArgs[0]] = labelArgs[1]
				}
			}
			result.Sources = append(result.Sources, s)
		}
	}
	return result, nil
}
