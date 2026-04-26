package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Resp 定义响应结构体形式
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

//code——————200 = 成功，400 = 参数错误，500 = 服务器错误
//msg————提示信息
//data————实际的数据内容

// SuccessResp 成功时返回数据给前端
func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// FailResp 操作不符合规定时返回操作提示
func FailResp(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Resp{
		Code: 400,
		Msg:  msg,
	})
}

// Page 分页，适用于数据量大的时候
//func Page(c *gin.Context, total int, list interface{}) {
//	c.JSON(http.StatusOK, gin.H{
//		"code":  200,
//		"msg":   "success",
//		"total": total,
//		"list":  list,
//	})
//}

// ErrorResp 服务器出错的时候返回
//func ErrorResp(c *gin.Context, msg string) {
//	c.JSON(http.StatusInternalServerError, Resp{
//		Code: 500,
//		Msg:  msg,
//	})
//}
