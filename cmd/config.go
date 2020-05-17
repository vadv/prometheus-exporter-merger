package cmd

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type source struct {
	Url    string            `yaml:"url"`
	Labels map[string]string `yaml:"labels"`
}

type config struct {
	Listen  string        `yaml:"listen"`
	Timeout time.Duration `yaml:"timeout"`
	Sources []*source     `yaml:"sources"`
}

func parseConfig(filename string) (*config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result := &config{
		Listen:  ":8080",
		Timeout: 15 * time.Second,
	}
	return result, yaml.Unmarshal(data, result)
}
