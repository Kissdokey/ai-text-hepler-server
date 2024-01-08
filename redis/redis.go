package redis

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// 创建redis客户端
var ctx = context.Background()
var rdb *redis.Client

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis地址，redis默认6379
		Password: "OkeyDokey",      // 自己本机redis密码，安装redis，修改密码，开机启动网上有教程
		DB:       0,                // 使用默认数据库
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	rdb = client
	return client
}
func AddUser(key string, value string) error {
	err := rdb.Set(ctx, key, value, 0).Err()
	return err
}
func HasUser(key string) error {
	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	if n <= 0 {
		return nil
	}else {
		return errors.New("账号重复")
	}
}
func Vertify(key string, value string) error {
	n, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if n > 0 {
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return err
		}
		if val == value {
			fmt.Printf("正确")
			return nil
		}
		fmt.Printf("value不对")
		return errors.New("存在key但是value不对")
	} else {
		fmt.Printf("不存在key:%s", key)
		return errors.New("不存在key")
	}
}
