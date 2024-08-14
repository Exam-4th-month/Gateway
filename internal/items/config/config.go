package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server ServerConfig
		Redis  RedisConfig
		JWT    JWTConfig
		Kafka  KafkaConfig
	}
	JWTConfig struct {
		SecretKey string
	}
	ServerConfig struct {
		ServerPort    string
		AuthPort      string
		BudgetingPort string
	}
	RedisConfig struct {
		Host string
		Port string
	}
	KafkaConfig struct {
		Brokers string
	}
)

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.Server.ServerPort = ":" + os.Getenv("SERVER_PORT")
	c.Server.AuthPort = ":" + os.Getenv("AUTH_PORT")
	c.Server.BudgetingPort = ":" + os.Getenv("BUDGETING_PORT")
	c.Redis.Host = os.Getenv("REDIS_HOST")
	c.Redis.Port = os.Getenv("REDIS_PORT")
	c.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	c.Kafka.Brokers = os.Getenv("KAFKA_BROKER_URI")

	return nil
}

func New() (*Config, error) {
	var config Config
	if err := config.Load(); err != nil {
		return nil, err
	}
	return &config, nil
}
