package api

import (
	"go-forum/dao"
	"go-forum/service"
	"go-forum/utils"

	"github.com/gin-gonic/gin"
)

// UserReq 注册/登录请求参数
type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailResp(c, "参数错误")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.FailResp(c, "用户名或密码不能为空！")
		return
	}
	err := service.SetUser(req.Username, req.Password)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, "注册成功！")
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var req UserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailResp(c, "参数错误")
		return
	}
	token, err := service.Login(req.Username, req.Password)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{"token": token})
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		utils.FailResp(c, "用户未登录!")
		return
	}
	userId := id.(uint)
	user, err := dao.GetUserById(userId)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{
		"id":       user.Id,
		"username": user.Username,
	})
}
