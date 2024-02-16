package middleware

import (
	"ai-text-helper-server/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ArnMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		userInfo, err := mysql.GetUserInfo(username.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "获取请求次数资源失败，请联系开发人员!",
			})
			c.Abort()
			return
		}
		if userInfo.ApiRequestNumber <= 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "剩余请求次数不足，请充值或者激活！",
			})
			c.Abort()
			return
		}
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
