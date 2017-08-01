package model

import (
	"log"
	"testing"
)

//var user = NewUser("root", "123456")

//var loginUser = NewLoginUser(10003, "123456")

////func TestCreateUser(t *testing.T) {
////	uid := user.CreateUser()
////	println(uid)
////}

//func TestAddFriend(t *testing.T) {
//	loginUser.AddFriend(10004)
//}

func TestAddMessage(t *testing.T) {
	AddMessage(1, Message{2, "hello"})
	AddMessage(1, Message{3, "hello"})
	log.Println(string(GetAllMessage(1)))
}


registeredUser := NewUser(10000)	//构造一个已经注册过的用户
unregisterUser := NewUser(10001)	//构造一个尚未注册的用户
func TestAddFriend(t *testing.T){
		//测试已注册但未登陆的用户进行好友添加
		uid := registeredUser.AddFriend(10010)
		if uid != 0 {
			log.Println("已注册未登录，添加好友成功") 
		}else {
			log.Println("已注册未登录，添加好友失败") 
		}
		//测试已注册并登陆的用户进行好友添加
		registeredUser.Login()
		uid = registeredUser.AddFriend(10010)
		if uid != 0 {
			log.Println("已注册已登录，添加好友成功") 
		}else {
			log.Println("已注册已登录，添加好友失败") 
		}
		//测试未注册未登录的用户进行好友添加
		uid = unregisterUser.AddFriend(10010)
			uid := registeredUser.AddFriend(10010)
		if uid != 0 {
			log.Println("未注册未登录，添加好友成功") 
		}else {
			log.Println("未注册未登录，添加好友失败") 
		}
}
