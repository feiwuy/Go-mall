// Package routers routes the router
package routers

import (
	"gomall/controllers"
	"gomall/middleware"

	// "net/http"

	"github.com/gin-gonic/gin"
)

// SetUpRouter 创建router
func SetUpRouter() *gin.Engine {
	r := gin.Default()

	// AI写的测试页面：由后端同源托管，避免浏览器 CORS 问题
	// 打开 http://localhost:3001/ 即可
	r.StaticFile("/", "./web/index.html")

	user := r.Group("/api/user")
	{
		user.POST("register", controllers.Register)
		user.POST("login", controllers.Login)
	}

	user.Use(middleware.AuthMiddleware())
	{
		user.GET("profile", controllers.GetProfile)
		user.PATCH("updateProfile", controllers.UpdateProfile)
	}

	todo := r.Group("api/todo")
	todo.Use(middleware.AuthMiddleware())
	{
		todo.POST("create", controllers.CreateTodo)
		todo.GET("getTodos", controllers.GetTodos)
		todo.GET("getTodo/:id", controllers.GetTodoByID)
		todo.PATCH("update/:id", controllers.UpdateTodo)
		todo.DELETE("delTodo/:id", controllers.DelTodo)
	}

	return r
}
