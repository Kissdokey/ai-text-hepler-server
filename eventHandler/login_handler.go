package handler

import (
	"ai-text-helper-server/mysql"
	"ai-text-helper-server/redis"
	"ai-text-helper-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Binding from JSON
type Login struct {
	User     string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// 读取密码账号，redis中比对，相等则生成jwt返回，前端和拿到jwt保存住
func LoginHandler(c *gin.Context) {
	var login Login
	if err := c.BindJSON(&login); err == nil {
		//Todo,密码加密
		if vertify := redis.Vertify(login.User, utils.MD5V(login.Password)); vertify != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  vertify.Error(),
				"data": gin.H{},
			})
		} else {
			//生成jwt返回
			jwt, jerr := utils.MakeToken(login.User, utils.MD5V(login.Password))
			if jerr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "生成token失败，请联系后端",
					"data": gin.H{},
				})
			} else {
				data, err := mysql.GetUserInfo(login.User)
				if err != nil {
					data.Avatar = ""
					data.ApiRequestNumber = 0
				}
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "登录成功",
					"data": gin.H{"jwt": jwt, "username": login.User, "avatar": data.Avatar, "apiRequestNumber": data.ApiRequestNumber},
				})
			}
		}
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"code": 406,
			"msg":  "参数错误，请按照文档请求参数",
			"data": gin.H{},
		})
	}
}
