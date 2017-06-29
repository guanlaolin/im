package db

import (
	"database/sql"
	"log"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

func GetDBConn() *sql.DB {
	db, err := sql.Open("mysql", "root:root@/db_im")
	if err != nil || db == nil {
		log.Println("连接数据库失败")
		return nil
	}
	return db
}

func GetNoConn() redis.Conn {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println("连接Redis失败", err.Error())
		return nil
	}
	return conn
}
