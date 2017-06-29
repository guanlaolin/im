package controller

import (
	"encoding/json"
	"log"
	"model"

	"github.com/gorilla/websocket"
)

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
