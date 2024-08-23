package api

import (
	_ "auth/api/docs"
	"auth/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router
// @title User
// @version 1.0
// @description API Gateway of Authorazation
// @host localhost:8085
// BasePath: /
func Router(hand *handler.Handler) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/register", hand.Register)
		auth.POST("/login", hand.Login)
		auth.POST("/change/password", hand.UpdatePassword)
		auth.POST("/refresh", hand.Refresh)
	}
	return router
}
