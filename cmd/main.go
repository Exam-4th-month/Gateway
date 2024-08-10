package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	casbin "github.com/casbin/casbin/v2"
	amqp "github.com/rabbitmq/amqp091-go"

	"gateway-service/internal/items/config"
	"gateway-service/internal/items/http/app"
	"gateway-service/internal/items/http/handler"
	"gateway-service/internal/items/redisservice"
	redisCl "gateway-service/internal/pkg/redis"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	redis, err := redisCl.NewRedisDB(config)
	if err != nil {
		log.Fatal(err)
	}

	modelPath := filepath.Join("internal", "items", "casbin", "model.conf")
	policyPath := filepath.Join("internal", "items", "casbin", "policy.csv")

	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := amqp.Dial(config.RabbitMQ.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	handler := handler.New(redisservice.New(redis, logger), logger, config, ch)

	log.Fatal(app.Run(handler, logger, config, enforcer))
}
