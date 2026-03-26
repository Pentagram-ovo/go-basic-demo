package main

import (
	"fmt"

	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/config"
	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/dao"
	"github.com/Pentagram-ovo/go-basic-demo/gorm-mysql-demo/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	config.InitMysql()
	user := &model.User{Username: "汪仔", Password: "2005"}
	err := dao.CreateUser(user)
	if err != nil {
		fmt.Println("创建失败！原因是：", err)
	}
	post := &model.Post{Title: "五角星发的第一个完美帖子！", Content: "这是完美帖子的内容~", UserId: 1}
	dao.CreatePostByTx(post)
}
