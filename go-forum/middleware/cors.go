package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 可以使得前端跨域调用后端接口
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		//开发环境：用 * 允许所有人跨域
		//生产环境：改成前端域名 https://xxx.com
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理 OPTIONS 预检
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
