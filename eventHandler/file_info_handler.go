package handler

import (
	"ai-text-helper-server/filePermission"
	"ai-text-helper-server/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateFile(c *gin.Context) {
	var data mysql.SQLData
	username, _ := c.Get("username")
	if err := c.BindJSON(&data); err == nil {
		canIEdit, editErr := filePermission.CanIEdit(username.(string), data.FileId)

		if editErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，文件更新失败111",
				"data": gin.H{},
			})
			return
		}
		if !canIEdit {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "当前文档您没有编辑权限，请联系作者申请权限",
				"data": gin.H{},
			})
			return
		}
		err := mysql.UpdateFile(data, data.FileId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，文件更新失败",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"code": 200,
				"msg":  "文件更新成功",
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

func CreateFile(c *gin.Context) {
	var data mysql.SQLData
	if err := c.BindJSON(&data); err == nil {
		err := mysql.UpdateFile(data, data.FileId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，文件更新失败",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"code": 200,
				"msg":  "文件更新成功",
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

func AcquireFile(c *gin.Context) {
	var data mysql.SQLData
	username, _ := c.Get("username")
	if err := c.BindJSON(&data); err == nil {
		canIVisible, visibleErr := filePermission.CanISee(username.(string), data.FileId)
		myPermission := filePermission.GetPermission(username.(string), data.FileId)
		if visibleErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，文件获取失败",
				"data": gin.H{},
			})
			return
		}
		if !canIVisible {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "当前文档您没有查看权限",
				"data": gin.H{},
			})
			return
		}
		file, err := mysql.GetFileInfo(data.FileId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 501,
				"msg":  "内部错误，文件获取失败",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"code": 200,
				"msg":  "文件获取成功",
				"data": gin.H{"author": file.Author,
					"fileName":         file.FileName,
					"content":          file.Content,
					"comments":         file.Comments,
					"lastModifiedTime": file.LastModifiedTime,
					"createTime":       file.CreateTime,
					"permissions":      file.Permissions,
					"selfPermission":   myPermission,
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
