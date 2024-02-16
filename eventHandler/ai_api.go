package handler

import (
	// "ai-text-helper-server/redis"
	//  "ai-text-helper-server/utils"
	"ai-text-helper-server/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Translate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": gin.H{},
	})
}
func Polish(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": gin.H{},
	})
}
func Authentication(c *gin.Context) {
	//能到达这里说明权限无误，可以登录，返回用户信息
	username, _ := c.Get("username")
	data,err := mysql.GetUserInfo(username.(string))
	if err != nil {
		data.Avatar = ""
		data.ApiRequestNumber = 0
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"username": username,
			"avatar": data.Avatar, 
			"apiRequestNumber": data.ApiRequestNumber,
		},
	})
}
