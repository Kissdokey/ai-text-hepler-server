package main

import "github.com/gin-gonic/gin"
import "ai-text-helper-server/routers"
import "ai-text-helper-server/redis"
import "ai-text-helper-server/utils"

func main() {
	// 创建路由
	r := gin.Default()
	// 注册路径
	routers.UserAccount(r)
	//初始化redis
	 redis.InitRedis()
	utils.GetEnvSalt()
	utils.GetEnvJwt()
	r.Run(":8000")
}
