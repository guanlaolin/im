package main

import (
	"controller"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //给日志加上文件名和行数
}

func main() {
	//go Ping()
	http.HandleFunc("/", controller.IndexHandle)
	http.HandleFunc("/ws", controller.WSHandle)               //聊天主页
	http.HandleFunc("/login", controller.LoginHandle)         //用户登录
	http.HandleFunc("/logout", controller.LogoutHandle)       //用户注销
	http.HandleFunc("/register", controller.RegisterHandle)   //用户注册
	http.HandleFunc("/search", controller.SearchUserHandle)   //查找用户
	http.HandleFunc("/userinfo", controller.UserInfoHandle)   //获取用户信息
	http.HandleFunc("/addfriend", controller.AddFriendHandle) //添加好友
	http.HandleFunc("/delete", controller.DelFriendHandle)    //删除好友 有bug，删除好友没有顺带删除未读消息
	http.HandleFunc("/unread", controller.UnReadMsgHandle)    //未读消息
	http.HandleFunc("/updatepsw", controller.UpdatePswHandle) //重置密码
	http.HandleFunc("/success", controller.SuccessHandle)     //成功

	//静态文件处理
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("../static")))) //静态文件
	http.Handle("/upload/",
		http.StripPrefix("/upload/",
			http.FileServer(http.Dir("../upload")))) //下载文件
	//err := http.ListenAndServe(":8888", nil) //监听8888端口
	//SSL
	err := http.ListenAndServeTLS(":8888", "../etc/cacert.pem", "../etc/privtkey.pem", nil)
	if err != nil {
		log.Fatalln("Listen error", err.Error())
	}
}

//心跳包
func Ping() {
	for {
		time.Sleep(5 * time.Second)
		for k, _ := range controller.Conns {
			err := controller.Conns[k].WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				delete(controller.Conns, k)
				log.Println(k, "已经从服务器断开")
				continue
			}
			log.Println("往", k, "发送心跳包成功")
		}
	}
}
