package routers


import (
	"ai-text-helper-server/eventHandler"
	"ai-text-helper-server/middleware"
	"github.com/gin-gonic/gin"
)

func FileInfo(r *gin.Engine) {
	fileInfo := r.Group("fileInfo")
	fileInfo.Use(middleware.JWTAuthMiddleware())
	{
		fileInfo.POST("/createFile",handler.CreateFile)
		fileInfo.POST("/updateFile",handler.UpdateFile)
		fileInfo.POST("/acquireFile",handler.AcquireFile)
	}
}

