package app

import (
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

//nolint
var (
	BuildCommit string
	BuildTime   string
)

type Config struct {
	QueryTimeout time.Duration `yaml:"query_timeout"`
}

func NewConfig(file io.Reader) (Config, error) {
	c := Config{}
	if err := yaml.NewDecoder(file).Decode(&c); err != nil {
		return Config{}, fmt.Errorf("error when decode app config: %w", err)
	}
	return c, nil
}
