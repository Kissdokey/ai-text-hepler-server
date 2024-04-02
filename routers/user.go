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
		userInfo.GET("/authentication", handler.Authentication)
		userInfo.POST("/updateAvatar", handler.UpdateAvatar)
		userInfo.POST("/updateDirectoryDependence",handler.UpdateDirectoryDependence)
		userInfo.POST("/updateUserFile",handler.UpdateUserFile)
		userInfo.POST("/updateFile",handler.UpdateFile)
	}
}

