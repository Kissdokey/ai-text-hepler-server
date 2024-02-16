package routers

import (
	"ai-text-helper-server/middleware"
	"github.com/gin-gonic/gin"
	"ai-text-helper-server/eventHandler"
)

func AiRequest(r *gin.Engine) {
	aiTextHelper := r.Group("aiTextHelper")
	aiTextHelper.Use(middleware.JWTAuthMiddleware())
	aiTextHelper.Use(middleware.ArnMiddleware())
	{
		aiTextHelper.POST("/textDeal", handler.TextDeal)
		aiTextHelper.POST("/customize", handler.Customize)
	}
}
