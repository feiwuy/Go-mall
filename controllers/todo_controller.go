package controllers

import (
	"errors"
	"gomall/global"
	"gomall/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateTodo 给用户创建一个todo
func CreateTodo(ctx *gin.Context) {
	var todo models.TodoList

	if err := ctx.ShouldBindJSON(&todo); err != nil {
		log.Printf("failed to create TodoList, error = %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法请求"})
		return
	}

	// 设置todo的ID(在中间件里定义过)
	todo.UserID = ctx.GetUint("user_id")

	// 给对应用户新建todo
	if err := global.DB.Create(&todo).Error; err != nil {
		log.Printf("failed to create a todoList, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg":"创建成功", "Title":todo.Title, "Content":todo.Content, "Completed":todo.Completed, "id": todo.ID})
}

func GetTodos(ctx *gin.Context) {
	var todos []models.TodoList

	if err := global.DB.Where("user_id = ?", ctx.GetUint("user_id")).Find(&todos).Error; err != nil {
		log.Printf("failed to get todos, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"todos": todos})
}

func GetTodoByID(ctx *gin.Context) {
	var todo models.TodoList
	id := ctx.Param("id")

	if err := global.DB.Where("user_id = ? AND id = ?", ctx.GetUint("user_id"), id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "不存在或无权访问"})
			return
		}
		log.Printf("failed to get todo, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "获取成功", "Title": todo.Title, "Content": todo.Content, "Completed": todo.Completed, "id": todo.ID})
}

func UpdateTodo(ctx *gin.Context) {
	type Update struct {
		Title		string	`json:"title"`
		Context   	string 	`json:"content"`
		Completed 	bool   	`json:"completed"`
	}
	var update Update

	id := ctx.Param("id")

	var todo models.TodoList
	if err := global.DB.Where("user_id = ? AND id = ?", ctx.GetUint("user_id"), id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "不存在或无权访问"})
			return
		}
		log.Printf("failed to get todo, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	if err := ctx.ShouldBindJSON(&update); err != nil {
		log.Printf("failed to update todo, error = %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法请求"})
		return
	}

	todo.Title = update.Title
	todo.Content = update.Context
	todo.Completed = update.Completed

	if err := global.DB.Save(&todo).Error; err != nil {
		log.Printf("failed to update todo, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
}

func DelTodo(ctx *gin.Context) {
	var todo models.TodoList
	id := ctx.Param("id")

	if err := global.DB.Where("user_id = ? AND id = ?", ctx.GetUint("user_id"), id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "不存在或无权访问"})
			return
		}
		log.Printf("failed to del todo, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	if err := global.DB.Delete(&todo).Error; err != nil {
		log.Printf("failed to del todo, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
