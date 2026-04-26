package api

import (
	"go-forum/service"
	"go-forum/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentReq 请求参数
type CommentReq struct {
	PostId  uint   `json:"post_id"`
	Content string `json:"content"`
}

// CommentSet 发布评论
func CommentSet(c *gin.Context) {
	var req CommentReq
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
	err := service.CreateComment(userid, req.PostId, req.Content)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, "发布评论成功！")
}

// CommentsListByPostId 根据postid查询列表接口
func CommentsListByPostId(c *gin.Context) {
	// 获取页码，默认 1
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	// 获取每页条数，默认 10
	sizeStr := c.Query("size")
	size := 10
	if sizeStr != "" {
		size, _ = strconv.Atoi(sizeStr)
	}
	//获取帖子id
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.FailResp(c, "帖子ID格式错误")
		return
	}
	postID := uint(id)
	comments, total, err := service.GetCommentListByPostID(postID, page, size)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{
		"comments": comments,
		"total":    total,
	})
}

// CommentsByUserId 根据用户id查询其发布的所有评论
func CommentsByUserId(c *gin.Context) {
	id, exists := c.Get("userID")
	if !exists {
		utils.FailResp(c, "用户未登录!不能查看其发布的评论！")
		c.Abort()
		return
	}
	userid := id.(uint)
	comments, err := service.GetCommentByUserID(userid)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{
		"comments": comments,
	})
}
