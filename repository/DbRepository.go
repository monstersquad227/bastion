package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	MysqlClient *sql.DB
)

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		"root",
		"1qaz@WSX",
		"192.168.1.87",
		"3307",
		"devflow",
		"utf8")
	var err error
	MysqlClient, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("打开 MySQL 连接失败: %v", err)
	}
	if err = MysqlClient.Ping(); err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}
	log.Println("MySQL 连接成功")
}
