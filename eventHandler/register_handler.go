package handler

import "net/http"
import "github.com/gin-gonic/gin"
import "ai-text-helper-server/redis"
import "ai-text-helper-server/utils"

//Todo:添加账号密码安全性校验
func RegisterHandler(c *gin.Context) {
	var login Login
	if err := c.ShouldBind(&login); err == nil {
		//判断账号是否存在
		if isUserExist := redis.HasUser((login.User)); isUserExist != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": 400,
				"msg":    "账号已存在",
				"data":   gin.H{},
			})
		} else {
			err := redis.AddUser(login.User, utils.MD5V(login.Password))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": 500,
					"msg":    "注册失败，请联系后端",
					"data":   gin.H{},
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": 200,
					"msg":    "注册成功",
					"data":   gin.H{},
				})
			}
		}
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status": 406,
			"msg":    "参数错误，请按照文档请求参数",
			"data":   gin.H{},
		})
	}
}
