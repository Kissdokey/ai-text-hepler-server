package main

import "github.com/gin-contrib/cors"
import "github.com/gin-gonic/gin"
import "ai-text-helper-server/routers"
import "ai-text-helper-server/redis"
import "ai-text-helper-server/mysql"
import "time"

// import "ai-text-helper-server/eventHandler"
// import "ai-text-helper-server/utils"

func main() {
	// 创建路由
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowWildcard:    true,
		AllowOrigins:     []string{"*"}, //后面申请到域名需要修改
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authentication"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	// 注册路径
	routers.InitRouters(r)
	//初始化redis
	redis.InitRedis()
	mysql.InitMySQL()
	r.Run(":8000")
}
