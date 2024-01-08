package routers

import "github.com/gin-gonic/gin"
import "ai-text-helper-server/eventHandler"

func UserAccount(r *gin.Engine) {
	userAccount := r.Group("userAccount")
	{
		userAccount.POST("/login", handler.LoginHandler)
		userAccount.POST("/register", handler.RegisterHandler)
	}
}
