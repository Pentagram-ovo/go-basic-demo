package middleware

import (
	"go-forum/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 token
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			utils.FailResp(c, "请先登录")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		// 2. 解析 token
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.FailResp(c, "登录已过期")
			c.Abort()
			return
		}

		// 3. 把用户ID存入上下文
		c.Set("userID", claims.UserID)

		// 4. 继续执行接口
		c.Next()
	}
}
