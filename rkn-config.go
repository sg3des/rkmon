package main

import (
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	HTTPaddr string `yaml:"http-addr"`

	GitRepo         string        `yaml:"git-repo"`
	GitRepoFilename string        `yaml:"git-repo-filename"`
	GitInterval     time.Duration `yaml:"git-interval"`

	Debug bool

	SMTP SMTP

	Subscribers []*Subscriber
}

func NewConfig(filename string) (c *Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return c, err
	}

	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return c, err
	}

	return
}
