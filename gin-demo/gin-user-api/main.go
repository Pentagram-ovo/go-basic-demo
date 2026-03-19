package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// User 定义用户结构体，用于JSON参数绑定和响应返回
// json标签：映射JSON字段名；binding标签：参数校验规则（name必填，age≥0）
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required"` // required：姓名为必填项
	Age  int    `json:"age" binding:"gte=0"`     // gte=0：年龄必须大于等于0
}

// 创建用户列表，储存每个用户的信息
// 全局变量：模拟内存数据库存储用户信息（key=用户ID，value=用户信息）
var userList = make(map[int]User)

// 全局变量：记录最后一个用户ID，用于生成新用户的自增ID
var last_user_id = 0

// StatCost 自定义全局中间件：统计每个HTTP请求的耗时，并生成随机请求ID
// 作用：1. 记录请求开始时间；2. 生成10000以内的随机请求ID并存入上下文；3. 执行后续路由逻辑；4. 计算并打印请求耗时
func StatCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		request_id := rand.Intn(10000)  //随机生成用户的Id
		c.Set("request_id", request_id) // 将请求ID存入Gin上下文，供后续接口使用
		c.Next()                        // 执行后续的路由处理函数（核心：不调用则后续逻辑不执行）
		//c.Abort()也可以用这个来终结程序进行
		cost := time.Since(start)
		log.Printf("请求方法为:%s 请求路径为：%s 请求Id为：%d 耗时为：%v", c.Request.Method, c.Request.URL.Path, request_id, cost)
	}
}

func main() {
	// gin.Default()：初始化Gin默认引擎，包含两个核心中间件：
	// 1. gin.Logger()：记录HTTP请求日志；2. gin.Recovery()：捕获panic并返回500错误，避免服务崩溃
	r := gin.Default()
	r.Use(StatCost())            //注册全局中间件
	user := r.Group("/api/user") //路由分组
	// GET /api/user/:id：根据路径参数id查询单个用户信息
	// 逻辑：1. 解析并校验ID格式；2. 查询内存数据库；3. 返回用户信息或错误提示
	user.GET("/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":      "参数不存在",
				"message":    "Id输入格式错误",
				"request_id": c.MustGet("request_id").(int),
			})
			return
		}
		user, ok := userList[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error":      "用户不存在",
				"message":    "该Id的用户还不存在",
				"request_id": c.MustGet("request_id").(int),
			})
			return
		}

		c.JSON(200, gin.H{
			"code":       200,
			"message":    "查询成功",
			"request_id": c.MustGet("request_id").(int),
			"data":       user,
		})
	})
	// POST /api/user：创建新用户（接收JSON参数）
	// 逻辑：1. 绑定并校验JSON参数；2. 生成自增用户ID；3. 存入内存数据库；4. 返回创建结果
	user.POST("", func(c *gin.Context) {
		var Newuser User
		// ShouldBindJSON：解析请求体中的JSON数据并绑定到Newuser，自动触发binding标签校验
		if err := c.ShouldBindJSON(&Newuser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      err.Error(),
				"message":    "参数校验失败",
				"request_id": c.MustGet("request_id").(int),
			})
			return
		}

		last_user_id++
		Newuser.Id = last_user_id
		userList[Newuser.Id] = Newuser

		c.JSON(http.StatusCreated, gin.H{
			"code":       201,
			"message":    "用户创建成功",
			"request_id": c.MustGet("request_id").(int),
			"data":       Newuser,
		})
	})
	// GET /api/user：多条件筛选用户（支持name/age查询参数，精准匹配）
	// 逻辑：1. 获取并解析查询参数；2. 遍历内存数据库筛选用户；3. 返回筛选结果（无匹配则提示）
	user.GET("", func(c *gin.Context) {
		name := c.Query("name")
		age_str := c.Query("age")
		age := -1
		if age_str != "" {
			ageInt, err := strconv.Atoi(age_str)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      err.Error(),
					"massage":    "年龄必须是大于0的整数!",
					"request_id": c.MustGet("request_id").(int),
				})
				return
			}
			age = ageInt
		}

		var result []User
		for _, u := range userList {
			if name != "" && u.Name != name {
				continue
			}
			if age != -1 && u.Age != age {
				continue
			}
			result = append(result, u)
		}
		len := len(result)
		if len == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":       400,
				"message":    "筛选查询失败！无符合的人员！",
				"request_id": c.MustGet("request_id").(int),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":       200,
				"message":    "筛选查询成功！",
				"request_id": c.MustGet("request_id").(int),
				"data":       result,
			})
		}
	})

	log.Println("服务启动成功，访问地址：http://localhost:8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("服务启动时失败：", err.Error())
	}
}
