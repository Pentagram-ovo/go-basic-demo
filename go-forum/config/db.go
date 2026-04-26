package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/go-forum?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}
	DB = db
	fmt.Println("连接数据库成功！")
}
