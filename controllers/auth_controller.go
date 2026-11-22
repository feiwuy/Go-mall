// Package controllers storages the program controller
package controllers

import (
	"gomall/global"
	"gomall/models"
	"gomall/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 存储注册逻辑
func Register(ctx *gin.Context) {
	var user models.User

	// 绑定数据给User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		log.Printf("failed to bind user, err = %v:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法请求"})
		return
	}

	if user.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}
	if len(user.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码至少6位"})
		return
	}

	// 判断用户名是否被注册
	var Count int64
	if err := global.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&Count).Error; err != nil {
		log.Printf("failed to register, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}
	if Count > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "该用户名已被占用"})
		return
	}

	// 将密码加密
	hashedPwd, err := utils.HashPwd(user.Password)
	if err != nil {
		log.Printf("failed to hash pwd, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	// 修改密码为加密后的
	user.Password = hashedPwd

	// 插入数据
	if err := global.DB.Create(&user).Error; err != nil {
		log.Printf("failed to insert data to database, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"msg": "注册成功!"})
}

// Login 存储登录逻辑
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to login, error = %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法请求"})
		return
	}

	var user models.User

	// 验证用户是否存在
	if err := global.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		log.Printf("failed to login, error = %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if !utils.CheckPassWord(input.Password, user.Password) {
		log.Printf("failed to login, error = password error")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := utils.GetToken(user.Username, user.ID)
	if err != nil {
		log.Printf("failed to generate token, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func GetProfile(ctx *gin.Context) {
	user := ctx.MustGet("user").(models.User)

	ctx.JSON(http.StatusOK, gin.H{"user": gin.H{"id": user.ID, "username": user.Username, "created_at": user.CreatedAt}})
}

func UpdateProfile(ctx *gin.Context) {
	var update struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&update); err != nil {
		log.Printf("failed to generate token, error = %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法请求"})
		return
	}

	if update.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	if len(update.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码至少6位"})
		return
	}

	userID := ctx.GetUint("user_id")

	// 检测用户是否被占用
	var count int64
	if err := global.DB.Model(&models.User{}).Where("username = ? AND id <> ?", update.Username, userID).Count(&count).Error; err != nil {
		log.Printf("failed to generate token, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "该用户名已被占用"})
		return
	}

	// 查找用户
	user := ctx.MustGet("user").(models.User)

	// 将密码加密
	hashedPwd, err := utils.HashPwd(update.Password)
	if err != nil {
		log.Printf("failed to hash pwd, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	// 更新用户信息
	user.Username = update.Username
	// 修改密码为加密后的
	user.Password = hashedPwd

	if err := global.DB.Save(&user).Error; err != nil {
		log.Printf("failed to update, error = %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误,请稍后再试"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": gin.H{"id": user.ID, "username": user.Username, "created_at": user.CreatedAt}})
}
