package model

import (
	"crypto/md5"
	"db"
	"fmt"
	"log"
)

//加密密码
func encrypt(password string) string {
	if password == "" {
		log.Println("密码不能为空，加密失败")
		return ""
	}
	b := md5.Sum([]byte(password))
	secret := fmt.Sprintf("%x", b)
	return secret
}

//测试用户合法性
func (user *User) checkUser() bool {
	if user.Uname == "" || user.Password == "" {
		log.Println("用户非法")
		return false
	}
	return true
}

//测试用户合法性
func (user *SummaryUser) checkSummaryUser() bool {
	if user.Uid == 0 || user.Uname == "" {
		log.Println("用户名非法")
		return false
	}
	return true
}

//保存到服务器的文件名
func GetFileKey(name string) int64 {
	if len(name) > 256 {
		log.Println("文件名长度大于256")
		return 0
	}
	conn := db.GetDBConn()
	if conn == nil {
		return 0
	}
	defer conn.Close()

	res, err := conn.Exec("INSERT INTO tb_key_file(file_name) VALUES(?)", name)
	if err != nil {
		log.Println("插入数据库tb_key_file失败", err)
		return 0
	}
	key, _ := res.LastInsertId()
	return key
}
