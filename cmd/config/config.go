package config

import (
	"github.com/rome314/idkb-events/internal/events"
	"github.com/rome314/idkb-events/pkg/connections"
)

var cfg *Config

type Config struct {
	ServerPort   string
	PgConnString string
	Redis        connections.RedisConfig
	App          events.Config
	PubSub       connections.RedisPubSubconfig
}

func GetConfig() *Config {
	return cfg
}
