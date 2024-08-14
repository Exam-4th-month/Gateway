package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	casbin "github.com/casbin/casbin/v2"
	"github.com/segmentio/kafka-go"

	"gateway-service/internal/items/config"
	"gateway-service/internal/items/http/app"
	"gateway-service/internal/items/http/handler"
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

	// time.Sleep(10 * time.Second)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.Kafka.Brokers},
		Logger:  log.New(os.Stdout, "kafka writer: ", 0),
	})
	defer writer.Close()

	handler := handler.New(logger, config, writer)

	log.Fatal(app.Run(handler, logger, config, enforcer))
}
