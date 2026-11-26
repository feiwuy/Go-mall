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
	// 打开 http://localhost:3000/ 即可
	r.StaticFile("/", "./web/index.html")

	user := r.Group("/api/users")
	{
		user.POST("", controllers.Register)
		user.POST("login", controllers.Login)
	}
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("", controllers.GetProfile)
		user.PATCH("", controllers.UpdateProfile)
	}

	todo := r.Group("/api/todos")
	todo.Use(middleware.AuthMiddleware())
	{
		todo.POST("", controllers.CreateTodo)
		todo.GET("", controllers.GetTodos)
		todo.GET("/:id", controllers.GetTodoByID)
		todo.PATCH("/:id", controllers.UpdateTodo)
		todo.DELETE("/:id", controllers.DelTodo)
	}

	return r
}
