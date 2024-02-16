package routers

import (
	"ai-text-helper-server/middleware"
	"github.com/gin-gonic/gin"
	"ai-text-helper-server/eventHandler"
)

func AiRequest(r *gin.Engine) {
	aiTextHelper := r.Group("aiTextHelper")
	aiTextHelper.Use(middleware.JWTAuthMiddleware())
	{
		aiTextHelper.GET("/authentication", handler.Authentication)
		aiTextHelper.POST("/translate", handler.Translate)
		aiTextHelper.POST("/polish", handler.Polish)
	}
}
