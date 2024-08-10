package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server   ServerConfig
		Redis    RedisConfig
		RabbitMQ RabbitMQConfig
	}
	ServerConfig struct {
		ServerPort  string
		AuthPort    string
		AthletePort string
		MedalPort   string
		EventPort   string
		StreamPort  string
	}
	RedisConfig struct {
		Host string
		Port string
	}
	RabbitMQConfig struct {
		RabbitMQ string
	}
)

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.Server.ServerPort = ":" + os.Getenv("SERVER_PORT")
	c.Server.AuthPort = ":" + os.Getenv("AUTH_PORT")
	c.Server.AthletePort = ":" + os.Getenv("ATHLETE_PORT")
	c.Server.MedalPort = ":" + os.Getenv("MEDAL_PORT")
	c.Server.EventPort = ":" + os.Getenv("EVENT_PORT")
	c.Server.StreamPort = ":" + os.Getenv("STREAM_PORT")
	c.Redis.Host = os.Getenv("REDIS_HOST")
	c.Redis.Port = os.Getenv("REDIS_PORT")
	c.RabbitMQ.RabbitMQ = os.Getenv("RABBITMQ_URI")

	return nil
}

func New() (*Config, error) {
	var config Config
	if err := config.Load(); err != nil {
		return nil, err
	}
	return &config, nil
}

// REDIS_URI=redis_uri
