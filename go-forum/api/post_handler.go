package api

import (
	"go-forum/service"
	"go-forum/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostReq 请求参数
type PostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// PostSet 发布帖子接口
func PostSet(c *gin.Context) {
	var req PostReq
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
	err := service.CreatePost(userid, req.Title, req.Content)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, "发布帖子成功！")
}

// PostGet 帖子详情接口
func PostGet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.FailResp(c, "帖子ID格式错误")
		return
	}
	postID := uint(id)
	// 调用带缓存的查询方法
	post, err := service.GetPostCache(postID)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{
		"id":      post.Id,
		"title":   post.Title,
		"content": post.Content,
		"userid":  post.UserID,
	})
}

// PostUpdate 修改更新帖子接口
func PostUpdate(c *gin.Context) {
	var req PostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailResp(c, "参数错误")
		return
	}
	//获取用户id
	useridTemp, exists := c.Get("userID")
	if !exists {
		utils.FailResp(c, "用户未登录!")
		return
	}
	userid := useridTemp.(uint)
	//要修改的帖子id在路径里输入
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.FailResp(c, "帖子ID格式错误")
		return
	}
	postID := uint(id)
	err = service.UpdatePost(postID, userid, req.Content, req.Title)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, "修改帖子成功！")
}

// PostDelete 删除帖子接口
func PostDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.FailResp(c, "帖子ID格式错误")
		return
	}
	postID := uint(id)
	err = service.DeletePost(postID)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, "删除帖子成功！")
}

// PostList 分页查询接口
func PostList(c *gin.Context) {
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
	list, total, err := service.GetPostListService(page, size)
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// GetHotPostList 获取热榜前n个帖子的详情
func GetHotPostList(c *gin.Context) {
	topStr := c.Query("top")
	top := 1
	if topStr != "" {
		top, _ = strconv.Atoi(topStr)
	}
	posts, err := service.GetTopNPosts(int64(top))
	if err != nil {
		utils.FailResp(c, err.Error())
		return
	}
	utils.SuccessResp(c, gin.H{
		"posts": posts,
	})
}
