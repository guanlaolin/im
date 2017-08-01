package main

import (
	"controller"
	"net/http"
)

//路由规则，暂时使用以下规则，后面修改为RESTful风格
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

//拟RESTful风格路由规则如下
//使用github.com/gorilla/mux包
//var urls = map[string]func(w http.ResponseWriter, r *http.Request){
//	"/":						//首页
//	"/ws":						//websocket通道
//	"/session/{uid}":			//会话信息，登录创建会话，注销销毁会话
//	"/user/{uid}":				//用户相关，如注册等
//	"/friend/{fid}":			//用户好友相关，如添加好友，删除好友等
//	"/message/{catagory}/{id}":	//消息
//}
