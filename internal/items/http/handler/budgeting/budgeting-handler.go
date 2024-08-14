package budgeting

import (
	"gateway-service/genproto/account"
	"gateway-service/genproto/budget"
	"gateway-service/genproto/category"
	"gateway-service/genproto/goal"
	"gateway-service/genproto/notification"
	"gateway-service/genproto/report"
	"gateway-service/genproto/transaction"
	"gateway-service/internal/items/config"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/items/redisservice"
	"log"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BudgetClientConn struct {
	AccountClient      account.AccountServiceClient
	BudgetClient       budget.BudgetServiceClient
	CategoryClient     category.CategoryServiceClient
	GoalClient         goal.GoalServiceClient
	NotificationClient notification.NotificationServiceClient
	ReportClient       report.ReportServiceClient
	TransactionClient  transaction.TransactionServiceClient
}

func NewBudgetClientConn(config *config.Config) *BudgetClientConn {
	return &BudgetClientConn{
		AccountClient:      account.NewAccountServiceClient(connect("budgeting", config.Server.BudgetingPort)),
		BudgetClient:       budget.NewBudgetServiceClient(connect("budgeting", config.Server.BudgetingPort)),
		CategoryClient:     category.NewCategoryServiceClient(connect("budgeting", config.Server.BudgetingPort)),
		GoalClient:         goal.NewGoalServiceClient(connect("budgeting", config.Server.BudgetingPort)),
		NotificationClient: notification.NewNotificationServiceClient(connect("budgeting", config.Server.BudgetingPort)),
		ReportClient:       report.NewReportServiceClient(connect("budgeting", config.Server.BudgetingPort)),
		TransactionClient:  transaction.NewTransactionServiceClient(connect("budgeting", config.Server.BudgetingPort)),
	}
}

type BudgetingHandler struct {
	AccountHandler      *AccountHandler
	BudgetHandler       *BudgetHandler
	CategoryHandler     *CategoryHandler
	GoalHandler         *GoalHandler
	NotificationHandler *NotificationHandler
	ReportHandler       *ReportHandler
	TransactionHandler  *TransactionHandler
}

func NewBudgetingHandler(redis *redisservice.RedisService, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *BudgetingHandler {
	clientConn := NewBudgetClientConn(config)

	return &BudgetingHandler{
		AccountHandler:      NewAccountHandler(redis, clientConn.AccountClient, logger, msgbroker, config),
		BudgetHandler:       NewBudgetHandler(clientConn.BudgetClient, logger, msgbroker, config),
		CategoryHandler:     NewCategoryHandler(clientConn.CategoryClient, logger, msgbroker, config),
		GoalHandler:         NewGoalHandler(clientConn.GoalClient, logger, msgbroker, config),
		NotificationHandler: NewNotificationHandler(clientConn.NotificationClient, logger, msgbroker, config),
		ReportHandler:       NewReportHandler(clientConn.ReportClient, logger, msgbroker, config),
		TransactionHandler:  NewTransactionHandler(clientConn.NotificationClient, clientConn.TransactionClient, logger, msgbroker, config),
	}
}

func connect(host, port string) *grpc.ClientConn {
	address := host + port
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
