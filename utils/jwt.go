package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type MyClaims struct {
	UserName             string `json:"userName"`
	PassWord             string `json:"password"`
	jwt.RegisteredClaims        // 注意!这是jwt-go的v4版本新增的，原先是jwt.StandardClaims
}

const MAXTIME = 24

func MakeToken(userName string, passWord string) (tokenString string, err error) {
	claim := MyClaims{
		UserName: userName,
		PassWord: passWord,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(MAXTIME * time.Hour * time.Duration(1))), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                             // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                             // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	jwtKey, _ := GetEnvJwt()
	MySecret := []byte(jwtKey) // 定义secret，后面会用到
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}
func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		key, _ := GetEnvJwt()
		return []byte(key), nil // 这是我的secret
	}
}

func ParseToken(tokenss string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenss, &MyClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
