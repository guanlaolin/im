package main

import (
	"controller"
	"net/http"
)

//路由规则
var urls = map[string]func(w http.ResponseWriter, r *http.Request){
	"/":          controller.IndexHandle,
	"/ws":        controller.WSHandle,         //聊天主页
	"/login":     controller.LoginHandle,      //用户登录
	"/logout":    controller.LogoutHandle,     //用户注销
	"/register":  controller.RegisterHandle,   //用户注册
	"/search":    controller.SearchUserHandle, //查找用户
	"/userinfo":  controller.UserInfoHandle,   //获取用户信息
	"/addfriend": controller.AddFriendHandle,  //添加好友
	"/delete":    controller.DelFriendHandle,  //删除好友 有bug，删除好友没有顺带删除未读消息
	"/unread":    controller.UnReadMsgHandle,  //未读消息
	"/updatepsw": controller.UpdatePswHandle,  //重置密码
	"/success":   controller.SuccessHandle,    //成功
}
