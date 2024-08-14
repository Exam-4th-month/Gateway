// Package api API.
//
// @title # Personal Finance Tracker
// @version 1.0
// @description API Endpoints for LocalEats
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package app

import (
	"log/slog"

	_ "gateway-service/internal/items/http/app/docs"
	"gateway-service/internal/items/middleware"

	"github.com/gin-contrib/cors"

	casbin "github.com/casbin/casbin/v2"

	"gateway-service/internal/items/config"
	"gateway-service/internal/items/http/handler"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(handler *handler.Handler, logger *slog.Logger, config *config.Config, enforcer *casbin.Enforcer) error {
	router := gin.Default()

	// CORS konfiguratsiyasi
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger dokumentatsiyasi uchun
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url, ginSwagger.PersistAuthorization(true)))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	superadmin := router.Group("superadmin")
	superadmin.Use(middleware.AuthzMiddleware("/superadmin", enforcer, config))
	{
		superadmin.POST("/createadmin", handler.AuthRepo.SuperAdminCreateAdminHandler)
	}

	auth := router.Group("auth")
	{
		superadmin := auth.Group("/superadmin")
		{
			superadmin.POST("/login", handler.AuthRepo.SuperAdminLoginHandler)
			superadmin.POST("/logout", handler.AuthRepo.SuperAdminLogoutHandler)
		}
		admin := auth.Group("/admin")
		{
			admin.POST("/login", handler.AuthRepo.AdminLoginHandler)
			admin.POST("/logout", handler.AuthRepo.AdminLogoutHandler)
		}
		user := auth.Group("/user")
		{
			user.POST("/register", handler.AuthRepo.RegisterHandler)
			user.POST("/login", handler.AuthRepo.LoginHandler)
			user.POST("/logout", handler.AuthRepo.LogoutHandler)
		}
	}

	admin := router.Group("admin")
	admin.Use(middleware.AuthzMiddleware("/admin", enforcer, config))
	{
		admin.PUT("/update/:id", handler.AuthRepo.UpdateUserHandler)
		admin.DELETE("/delete/:id", handler.AuthRepo.DeleteUserHandler)

	}

	user := router.Group("user")
	user.Use(middleware.AuthzMiddleware("/user", enforcer, config))
	{
		account := user.Group("account")
		{
			account.POST("/", handler.BudgetingRepo.AccountHandler.CreateAccountHandler)
			account.GET("/", handler.BudgetingRepo.AccountHandler.GetAccountsHandler)
			account.GET("/:id", handler.BudgetingRepo.AccountHandler.GetAccountByIdHandler)
			account.PUT("/", handler.BudgetingRepo.AccountHandler.UpdateAccountHandler)
			account.DELETE("/:id", handler.BudgetingRepo.AccountHandler.DeleteAccountHandler)
		}

		budget := user.Group("budget")
		{
			budget.POST("/", handler.BudgetingRepo.BudgetHandler.CreateBudgetHandler)
			budget.GET("/", handler.BudgetingRepo.BudgetHandler.GetBudgetsHandler)
			budget.GET("/:id", handler.BudgetingRepo.BudgetHandler.GetBudgetByIdHandler)
			budget.PUT("/", handler.BudgetingRepo.BudgetHandler.UpdateBudgetHandler)
			budget.DELETE("/:id", handler.BudgetingRepo.BudgetHandler.DeleteBudgetHandler)
		}

		category := user.Group("category")
		{
			category.POST("/", handler.BudgetingRepo.CategoryHandler.CreateCategoryHandler)
			category.GET("/", handler.BudgetingRepo.CategoryHandler.GetCategoriesHandler)
			category.GET("/:id", handler.BudgetingRepo.CategoryHandler.GetCategoryByIdHandler)
			category.PUT("/", handler.BudgetingRepo.CategoryHandler.UpdateCategoryHandler)
			category.DELETE("/:id", handler.BudgetingRepo.CategoryHandler.DeleteCategoryHandler)
		}

		goal := user.Group("goal")
		{
			goal.POST("/", handler.BudgetingRepo.GoalHandler.CreateGoalHandler)
			goal.GET("/", handler.BudgetingRepo.GoalHandler.GetGoalsHandler)
			goal.GET("/:id", handler.BudgetingRepo.GoalHandler.GetGoalByIdHandler)
			goal.PUT("/", handler.BudgetingRepo.GoalHandler.UpdateGoalHandler)
			goal.DELETE("/:id", handler.BudgetingRepo.GoalHandler.DeleteGoalHandler)
		}

		transaction := user.Group("transaction")
		{
			transaction.POST("/", handler.BudgetingRepo.TransactionHandler.CreateTransactionHandler)
			transaction.GET("/", handler.BudgetingRepo.TransactionHandler.GetTransactionsHandler)
			transaction.GET("/:id", handler.BudgetingRepo.TransactionHandler.GetTransactionByIdHandler)
			transaction.PUT("/", handler.BudgetingRepo.TransactionHandler.UpdateTransactionHandler)
			transaction.DELETE("/:id", handler.BudgetingRepo.TransactionHandler.DeleteTransactionHandler)
		}

		report := user.Group("report")
		{
			report.POST("/spending", handler.BudgetingRepo.ReportHandler.GetSpendingReportHandler)
			report.POST("/incoming", handler.BudgetingRepo.ReportHandler.GetIncomeReportHandler)
			report.POST("/bugdet", handler.BudgetingRepo.ReportHandler.GetBudgetPerformanceReportHandler)
			report.POST("/goal", handler.BudgetingRepo.ReportHandler.GetGoalProgressReportHandler)
		}

		notification := user.Group("notification")
		{
			notification.GET("/", handler.BudgetingRepo.NotificationHandler.GetNotifications)
			notification.PUT("/", handler.BudgetingRepo.NotificationHandler.MarkNotificationAsRead)
		}
	}

	return router.Run("gateway"+config.Server.ServerPort)
}
