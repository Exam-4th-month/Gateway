package handler

import (
	"log"
	"log/slog"

	auth_pb "gateway-service/genproto/auth"
	"gateway-service/internal/items/config"
	auth_broker "gateway-service/internal/items/msgbroker/auth"

	"gateway-service/internal/items/http/handler/auth"
	"gateway-service/internal/items/redisservice"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	Handler struct {
		AuthRepo *auth.AuthHandler
	}
)

func connect(port string) *grpc.ClientConn {
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func New(redis *redisservice.RedisService, logger *slog.Logger, config *config.Config, channel *amqp.Channel) *Handler {
	authClient := auth.NewAthleteHandler(logger, auth_pb.NewAuthServiceClient(connect(config.Server.AuthPort)), redis, auth_broker.NewAthleteMsgBroker(channel, logger))

	return &Handler{
		AuthRepo: authClient,
	}
}
