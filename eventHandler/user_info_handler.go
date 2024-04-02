package handler

import (
	"ai-text-helper-server/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Todo:添加账号密码安全性校验
func UpdateAvatar(c *gin.Context) {
	var data mysql.SQLData
	if err := c.BindJSON(&data); err == nil {
		username, _ := c.Get("username")
		err := mysql.UpdateAvatar(data.Avatar, username.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，头像未更新成功",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"code": 200,
				"msg":  "头像更新成功",
				"data": gin.H{
					"avatar": data.Avatar,
				},
			})
		}

	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"code": 406,
			"msg":  "参数错误，请按照文档请求参数",
			"data": gin.H{},
		})
	}
}
func UpdateDirectoryDependence(c *gin.Context) {
	var data mysql.SQLData
	if err := c.BindJSON(&data); err == nil {
		username, _ := c.Get("username")
		err := mysql.UpdateDirectoryDependence(username.(string), data.DirectoryDependence)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，目录依赖未更新成功",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"code": 200,
				"msg":  "目录依赖更新成功",
				"data": gin.H{},
			})
		}

	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"code": 406,
			"msg":  "参数错误，请按照文档请求参数",
			"data": gin.H{},
		})
	}
}

func UpdateUserFile(c *gin.Context) {
	var data mysql.SQLData
	if err := c.BindJSON(&data); err == nil {
		username, _ := c.Get("username")
		err := mysql.UpdateUserFile(username.(string), data.FileId, data.IsDelete)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，用户文件更新失败",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"code": 200,
				"msg":  "用户文件更新成功",
				"data": gin.H{},
			})
		}

	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"code": 406,
			"msg":  "参数错误，请按照文档请求参数",
			"data": gin.H{},
		})
	}
}


