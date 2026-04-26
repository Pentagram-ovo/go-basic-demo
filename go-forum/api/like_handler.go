package api

import (
	"go-forum/service"
	"go-forum/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Like struct {
	PostID uint `json:"post_id"`
}

// LikeAction 点赞/取消点赞操作
func LikeAction(c *gin.Context) {
	var req Like
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailResp(c, "参数错误")
		return
	}
	id, exists := c.Get("userID")
	if !exists {
		utils.FailResp(c, "用户未登录!")
		return
	}
	userid := id.(uint)
	err := service.ToggleLikePost(userid, req.PostID)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, "修改点赞成功！")
}

// GetLikeStatus 获取点赞状态
func GetLikeStatus(c *gin.Context) {
	idStr := c.Param("id")
	idTemp, err := strconv.Atoi(idStr)
	if err != nil {
		utils.FailResp(c, "帖子ID格式错误")
		return
	}
	postID := uint(idTemp)
	id, exists := c.Get("userID")
	if !exists {
		utils.FailResp(c, "用户未登录")
		return
	}
	userid := id.(uint)
	status, err := service.GetLikeStatus(userid, postID)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	if status {
		utils.SuccessResp(c, "用户已点赞！")
	} else {
		utils.SuccessResp(c, "用户未点赞！")
	}
}

// GetLikeCount 获取帖子点赞总数
func GetLikeCount(c *gin.Context) {
	idStr := c.Param("id")
	idTemp, err := strconv.Atoi(idStr)
	if err != nil {
		utils.FailResp(c, "帖子ID格式错误")
		return
	}
	postID := uint(idTemp)
	count, err := service.GetLikeCount(postID)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{"count": count})
}
