package utils

import (
	"errors"
	"fmt"
	"os"
)

func GetEnvSalt() (string, error) {
	// 获取密钥
	saltKey := os.Getenv("SALT_KEY")
	if saltKey == "" {
		fmt.Println("无法获取盐，请配置环境变量：SALT_KEY")
		return saltKey, errors.New("无法获取盐，请配置环境变量：SALT_KEY")
	} else {
		return saltKey, nil
	}
}
func GetEnvJwt() (string, error) {
	// 获取密钥
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		fmt.Println("无法获取jwt密钥，请配置环境变量：JWT_KEY")
		return jwtKey, errors.New("无法获取jwt密钥，请配置环境变量：JWT_KEY")
	} else {
		return jwtKey, nil
	}
}

func GetEnvSQL() (string, error) {
	// 获取密钥
	sqlPassword := os.Getenv("MYSQL_PASSWORD")
	if sqlPassword == "" {
		fmt.Println("无法获取MySQL密钥，请配置环境变量：MYSQL_PASSWORD")
		return sqlPassword, errors.New("无法获取MySQL密钥，请配置环境变量：MYSQL_PASSWORD")
	} else {
		return sqlPassword, nil
	}
}
func GetEnvApiKey() (string,error) {
	// 获取密钥
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("无法获取api-key，请配置环境变量：API_KEY")
		return apiKey, errors.New("无法获取api-key，请配置环境变量：API_KEY")
	} else {
		return apiKey, nil
	}
}
func GetEnvSecretKey() (string,error) {
	// 获取密钥
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("无法获取secret-key，请配置环境变量：SECRET_KEY")
		return secretKey, errors.New("无法获取secret-key，请配置环境变量：SECRET_KEY")
	} else {
		return secretKey, nil
	}
}
func GetEnvAccessToken() (string,error) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		fmt.Println("无法获取access-token，请配置环境变量：ACCESS_TOKEN")
		return accessToken, errors.New("无法获取secret-key，请配置环境变量：ACCESS_TOKEN")
	} else {
		return accessToken, nil
	}
}