package main

import (
	"flag"
	"logger"
	"net/http"
	"tool"

	"github.com/gorilla/mux"
)

var l *logger.Logger = logger.NewLoggerWithFile(logger.LEVEL_DEBUG, "../log/log")

func init() {
	flag.Parse()

	err := tool.ConfSetPath(*config)
	if err != nil {
		l.FATAL("初始化配置文件失败", err)
	}
}

func main() {
	//路由
	l.DEBUG("开始配置路由")
	router := mux.NewRouter()
	for url, handler := range urls {
		router.HandleFunc(url, handler)
	}
	http.Handle("/", router)
	l.DEBUG("路由配置完成")

	//静态文件处理
	//实际应用交由nginx解析，或者使用cdn，这里仅仅是为了开发方便
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("../static")))) //静态文件
	http.Handle("/upload/",
		http.StripPrefix("/upload/",
			http.FileServer(http.Dir("../upload")))) //下载文件

	//使用SSL
	l.DEBUG("启动监听端口", tool.Conf("addr"), "使用的证书是：", tool.Conf("cert"),
		"使用的key是：", tool.Conf("key"))
	err := http.ListenAndServeTLS(tool.Conf("addr"),
		tool.Conf("cert"), tool.Conf("key"), nil)
	if err != nil {
		l.FATAL("监听失败", err.Error())
	}
}
