package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/segmentio/kafka-go"

	"gateway-service/internal/items/config"
	"gateway-service/internal/items/http/app"
	"gateway-service/internal/items/http/handler"
	"gateway-service/internal/items/msgbroker"
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

	modelPath := filepath.Join("internal", "items", "casbin", "model.conf")
	policyPath := filepath.Join("internal", "items", "casbin", "policy.csv")

	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		log.Fatal(err)
	}

	redis, err := redisCl.NewRedisDB(config)
	if err != nil {
		logger.Error("Error connecting to Redis", slog.String("err", err.Error()))
	}

	time.Sleep(10 * time.Second)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.Kafka.Brokers},
		Logger:  log.New(os.Stdout, "kafka writer: ", 0),
	})
	defer writer.Close()

	err = msgbroker.CreateTopics(config, logger)
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.New(redisservice.New(redis, logger), logger, config, writer)

	log.Fatal(app.Run(handler, logger, config, enforcer))
}
