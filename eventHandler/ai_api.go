package handler

import (
	"ai-text-helper-server/ai_request"
	"ai-text-helper-server/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestStruct struct {
	RequestType    string `json:"requestType"`
	RequestContent string `json:"requestContent"`
}

func TextDeal(c *gin.Context) {
	var requestStruct RequestStruct
	if err := c.BindJSON(&requestStruct); err == nil {
		res, resErr := ai.AiRequest(requestStruct.RequestType, requestStruct.RequestContent)
		if resErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误或者参数错误！",
				"data": gin.H{},
			})
		} else {
			//请求次数减一，返回
			username, _ := c.Get("username")
			apiRequestNumber,err := mysql.UpdateRequestNumber(username.(string), -1)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 501,
					"msg":  "内部错误，请求次数更新失败",
					"data": gin.H{},
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "ok",
					"data": gin.H{
						"content":          res,
						"apiRequestNumber": apiRequestNumber,
					},
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
func Customize(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": gin.H{},
	})
}
func Authentication(c *gin.Context) {
	//能到达这里说明权限无误，可以登录，返回用户信息
	username, _ := c.Get("username")
	data, err := mysql.GetUserInfo(username.(string))
	if err != nil {
		data.Avatar = ""
		data.ApiRequestNumber = 0
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"username":         username,
			"avatar":           data.Avatar,
			"apiRequestNumber": data.ApiRequestNumber,
		},
	})
}
