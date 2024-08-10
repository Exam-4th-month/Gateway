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
	superadmin.Use(middleware.AuthzMiddleware("/superadmin", enforcer))
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
	admin.Use(middleware.AuthzMiddleware("/admin", enforcer))
	{
		auth := admin.Group("auth")
		{
			auth.PUT("/update/:id", handler.AuthRepo.UpdateUserHandler)
			auth.DELETE("/delete/:id", handler.AuthRepo.DeleteUserHandler)
		}

	}

	return router.Run(config.Server.ServerPort)
}
