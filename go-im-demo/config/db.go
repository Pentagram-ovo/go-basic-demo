package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 为docker做准备
func getDSN() string {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		password = "123456"
	}
	dbName := os.Getenv("MYSQL_DB")
	if dbName == "" {
		dbName = "im"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)
}

func WaitForMySQL() {
	for i := 0; i < 30; i++ {
		db, err := gorm.Open(mysql.Open(getDSN()), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			if sqlDB.Ping() == nil {
				sqlDB.Close()
				return
			}
			sqlDB.Close()
		}
		log.Printf("等待 MySQL 就绪... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("MySQL 未能就绪，退出")
}

func InitMysql() {
	dsn := getDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}
	DB = db
	fmt.Println("连接数据库成功！")
}
