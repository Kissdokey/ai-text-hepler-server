package routers

import "github.com/gin-gonic/gin"

func InitRouters(r *gin.Engine) {
	UserAccount(r)
	AiRequest(r)
	UserInfo(r)
	FileInfo(r)
}
