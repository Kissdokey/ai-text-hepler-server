package utils

import (
	"errors"
	"fmt"
	"os"
)  
  
func GetEnvSalt() (string,error) {  
	// 获取密钥  
	saltKey := os.Getenv("SALT_KEY")  
	if saltKey == "" {  
		fmt.Println("无法获取盐，请配置环境变量：SALT_KEY")  
		return saltKey,errors.New("无法获取盐，请配置环境变量：SALT_KEY")
	}else {
		return saltKey,nil
	}
}  
func GetEnvJwt() (string,error) {  
	// 获取密钥  
	jwtKey := os.Getenv("JWT_KEY")  
	fmt.Println(jwtKey)
	if jwtKey == "" {  
		fmt.Println("无法获取jwt密钥，请配置环境变量：JWT_KEY")  
		return jwtKey,errors.New("无法获取jwt密钥，请配置环境变量：JWT_KEY")
	}else {
		return jwtKey,nil
	}
}  