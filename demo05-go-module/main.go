package main

//第三方包gin 用go mod tidy下载
import (
	"demo05/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// ========== 1. 使用自定义包 ==========
	fmt.Println("===== 自定义包使用 =====")
	tools.SayHello()

	// ========== 2. 使用第三方包（gin） ==========
	fmt.Println("\n===== 第三方包（gin）使用 =====")
	// 创建 gin 实例
	r := gin.Default()
	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "五角星",
		})
	})

	// 启动 HTTP 服务（端口 8080）
	fmt.Println("Web 服务启动：http://localhost:8080")
	err := r.Run(":8080")
	if err != nil {
		fmt.Printf("服务启动失败：%v\n", err)
	}
}

//go build生成可执行的.exe文件
