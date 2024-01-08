### 配置redis密码
如果redis未安装，先安装，然后主机配置密码（安全性考虑，如不需要也可以）
./redis/resis.go中
func InitRedis
client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis地址，redis默认6379
		Password: "OkeyDokey",      // 自己本机redis密码，安装redis，修改密码，开机启动网上有教程
		DB:       0,                // 使用默认数据库
	})
这里password设置为自己本机的redis密码，如果没有密码就清空
### 配置系统变量
jwt密钥，请配置环境变量：JWT_KEY
md5加密盐，请配置环境变量：SALT_KEY
### 初始化 
`go mod tidy`
### 运行
`go run main.go`