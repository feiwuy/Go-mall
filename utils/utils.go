// Package utils 存储用户需要的工具
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// HashPwd 生成哈希后的密码并返回
func HashPwd(pwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hashed), err
}

// GetToken 获取JWTToken
func GetToken(username string) (string, error) {
	// 载荷里放过期时间和用户名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // 过期时间
	})
	// 用密钥进行签名, 签名不代表加密, 只是为了确保token是自己签发的
	singedToken, err := token.SignedString([]byte("secret"))
	return "Bearer " + singedToken, err
}

func CheckPassWord(Pwd, hashedPwd string) bool {
	// 先写真实密码, 需要比对的放右边
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(Pwd))
	return err == nil
}

// PraseJWT 解析JWT
func PraseJWT(tokenstring string) (string, error) {
	if len(tokenstring) >= 7 && tokenstring[:7] == "Bearer " {
		tokenstring = tokenstring[7:]
	}

	// 获取token
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		// 类型断言, 判断token的签名方法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected singing method!")
		}
		// 返回密钥
		return []byte("secret"), nil
	})

	if err != nil {
		return "", err
	}

	// 类型断言, 查看token的声明方法是否正确
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("Username claim is not a string")
		}
		return username, err
	}
	return "", err
}
