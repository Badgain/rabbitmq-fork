package config

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/Badgain/rabbit/config"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

type ForkConfig interface {
	GetMapping() []QueueExchangeMapping
	GetServerInfo() config.ServerConfig
}

type cfg struct {
	ServerConfig config.ServerConfig    `json:"server" yaml:"server"`
	Mapping      []QueueExchangeMapping `json:"mapping" yaml:"mapping"`
}

type QueueExchangeMapping struct {
	Queue       config.QueueConfig       `json:"queue" yaml:"queue"`
	Exchange    config.ExchangeConfig    `json:"exchange" yaml:"exchange"`
	MessageType config.MessageTypeConfig `json:"message_type" yaml:"message_type"`
}

func (c *cfg) GetMapping() []QueueExchangeMapping {
	return c.Mapping
}

func (c *cfg) GetServerInfo() config.ServerConfig {
	return c.ServerConfig
}

func NewConfig(lifecycle fx.Lifecycle) ForkConfig {
	c := &cfg{}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			path := os.Getenv("CONFIG_PATH")
			if path == "" {
				return errors.New("config path is empty")
			}

			ext := filepath.Ext(path)
			if ext == "" {
				return errors.New("config file extension is empty")
			}

			f, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(os.O_RDONLY))
			if err != nil {
				return err
			}

			defer f.Close()

			bts, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			switch ext {
			case ".json":
				if err = json.Unmarshal(bts, c); err != nil {
					return err
				}
			case ".yaml", ".yml":
				if err = yaml.Unmarshal(bts, c); err != nil {
					return err
				}
			default:
				return errors.New("unsupported config file extension")
			}

			return nil
		},
	})

	return c
}
