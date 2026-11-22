// Package middleware 存储中间件
package middleware

import (
	"gomall/global"
	"gomall/models"
	"gomall/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 是一个用来判断用户登陆的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求中获取Authorization
		token := ctx.GetHeader("Authorization")
		if token == "" {
			log.Print("failed to authorization, error = Missing Authorization Header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			// 停止接下来的任务
			ctx.Abort()
			return
		}

		// 解析token获取username
		username, userID, err := utils.ParseJWT(token)
		if err != nil {
			log.Printf("failed to prase token, error = %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// 设置上下文的username
		ctx.Set("username", username)

		// 判断用户是否被删除并获取id
		var user models.User
		if err := global.DB.First(&user, userID).Error; err != nil {
			log.Printf("failed to login, error = %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "登录失效,请重新登录"})
			ctx.Abort()
			return
		}

		// 将用户信息传入上下文
		ctx.Set("user_id", userID)
		ctx.Set("user", user)

		// 调用下一个处理函数
		ctx.Next()
	}
}
