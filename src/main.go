package main

import (
	"controller"
	"logger"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var l *logger.Logger = logger.NewLoggerWithFile(logger.LEVEL_DEBUG, "../log/mylog.log")

func init() {
}

func main() {
	//路由
	for url, handler := range urls {
		http.HandleFunc(url, handler)
	}

	//静态文件处理
	//实际应用交由nginx解析，这里仅仅是为了开发方便
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("../static")))) //静态文件
	http.Handle("/upload/",
		http.StripPrefix("/upload/",
			http.FileServer(http.Dir("../upload")))) //下载文件

	//SSL
	err := http.ListenAndServeTLS(":8888", "../etc/cacert.pem", "../etc/privtkey.pem", nil)
	if err != nil {
		l.FATAL("Listen error", err.Error())
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
				//log.Println(k, "已经从服务器断开")
				l.DEBUG(k, "已经从服务器断开")
				continue
			}
			l.DEBUG("往", k, "发送心跳包成功")
			//log.Println("往", k, "发送心跳包成功")
		}
	}
}
