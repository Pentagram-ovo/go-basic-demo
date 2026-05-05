package api

import (
	"go-im-demo/internal/middleware"
	"go-im-demo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 处理注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	user, err := h.UserService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "注册成功！",
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	user, err1 := h.UserService.Login(req.Username, req.Password)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err1.Error()})
		return
	}
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "登录成功"})
}
