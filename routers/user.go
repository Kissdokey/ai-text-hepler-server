package routers


import (
	"ai-text-helper-server/eventHandler"
	"ai-text-helper-server/middleware"
	"github.com/gin-gonic/gin"
)

func UserInfo(r *gin.Engine) {
	userInfo := r.Group("userInfo")
	userInfo.Use(middleware.JWTAuthMiddleware())
	{
		userInfo.POST("/updateAvatar", handler.UpdateAvatar)
	}
}

