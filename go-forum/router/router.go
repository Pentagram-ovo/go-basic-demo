package router

import (
	"go-forum/api"
	"go-forum/middleware"
	"go-forum/utils"

	"github.com/gin-gonic/gin"
)

// 路由总入口
func SetupRouter(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		utils.SuccessResp(c, "五角星的go-forum 论坛 API 服务运行中...")
	})

	publicGroup := r.Group("/api")
	{
		publicGroup.POST("/user/register", api.UserRegister)
		publicGroup.POST("/user/login", api.UserLogin)
		publicGroup.GET("/post/list", api.PostList)
		publicGroup.GET("/post/:id", api.PostGet)
		publicGroup.GET("/comment/:id", api.CommentsListByPostId)
		publicGroup.GET("/like/count/:id", api.GetLikeCount)
		publicGroup.GET("/post/hot", api.GetHotPostList)
	}

	privateGroup := r.Group("/api")
	{
		//私密用法先进行鉴权
		privateGroup.Use(middleware.Auth())
		{
			privateGroup.GET("/user/info", api.UserInfo)
			privateGroup.POST("/post/set", api.PostSet)
			privateGroup.PUT("/post/update/:id", api.PostUpdate)
			privateGroup.DELETE("/post/delete/:id", api.PostDelete)
			privateGroup.POST("/comment/set", api.CommentSet)
			privateGroup.GET("/comment/user", api.CommentsByUserId)
			privateGroup.POST("/like/action", api.LikeAction)
			privateGroup.GET("/like/status/:id", api.GetLikeStatus)
		}
	}
}
