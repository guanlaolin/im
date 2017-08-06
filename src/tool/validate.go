package tool

import "log"

//检查uid合法性
//uid >= 10000 || uid < ?
func ValidateUid(uid int) bool {
	if uid < 10000 {
		return false
	}
	return true
}

//检查密码合法性
//psw:
//	长度：6-20
func ValidatePsw(psw string) bool {
	if len(psw) < 6 || len(psw) > 20 {
		return false
	}
	return true
}

//检查用户名合法性
//uname
//	长度4-10
func ValidateUname(uname string) bool {
	if len(uname) < 4 || len(uname) > 10 {
		return false
	}
	return true
}

//检查邮箱合法性
func ValidateEmail(email string) bool {
	if len(email) < 6 || len(email) > 32 {
		log.Println("邮箱长度必须为6-32")
		return false
	}
	return true
}
