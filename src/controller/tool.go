package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"model"
	"net/http"
	"tool"

	"github.com/gorilla/websocket"
)

//解析并渲染模板文件
//已进行模板文件是否存在检测
func RenderTmpl(path string, data interface{}, w http.ResponseWriter) error {
	if err := tool.ExistFile(path); err != nil {
		l.DEBUG("模板文件不存在：", path)
		return err
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		l.DEBUG("解析模板文件错误:", path, err.Error())
		http.Error(w, "服务器错误，登录失败", http.StatusInternalServerError)
		return err
	}

	if err := tmpl.Execute(w, data); err != nil {
		l.DEBUG("渲染模板文件失败", err.Error())
	}
	return nil
}

//通过session判断是否已经登录，error为空，则已经登录
func IsLogin(r *http.Request) bool {
	session, err := store.Get(r, tool.Conf("session-name"))
	if err != nil {
		l.INFO("获取session失败")
		return false
	}
	if nil == session.Values["uid"] {
		return false
	}
	return true
}

func WSErrHandler(text string, ws *websocket.Conn) {
	log.Println("ws错误处理开始...")
	var msg model.Message
	msg.Type = model.ERROR
	msg.Value = text
	json_data, err := json.Marshal(msg)
	if err != nil {
		log.Println("WSErrHandler打包json失败", err)
	}
	ws.WriteMessage(model.ERROR, json_data)
	return
}
