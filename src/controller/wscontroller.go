//由于ws代码比较场，故专门分一个文件
package controller

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"model"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

//ws逻辑
//1、写入Conns
func WSHandle(w http.ResponseWriter, r *http.Request) {
	//记录当前用户
	var cur_uid int = 0

	//当前上传文件名
	var cur_file_name string

	//转化为websocket协议
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("Upgrade error")
		http.Error(w, "转websocket失败", http.StatusInternalServerError)
	}
	log.Println("websocket success")
	defer conn.Close()

	//循环读取报文
	for {
		var msg model.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			//bug
			log.Println("ReadJSON error", err)
		}

		log.Println("收到一条消息", msg)

		//消息格式选择
		switch msg.Type {
		case model.IDENTITY:
			//记录用户信息的报文
			//bug
			log.Println("判定为身份消息")
			cur_uid = msg.From
			if !model.CheckUid(cur_uid) {
				log.Println("用户id不合法，连接断开", err)
				conn.WriteMessage(websocket.CloseMessage, []byte("上送的uid有误"))
				return
			}
			Conns[cur_uid] = conn
			continue
		case model.TEXT:
			//文本消息

		case model.FILE:
			var key int64
			var file *os.File

			data_str := msg.Value.(string)
			//解base64，得到数据
			temp, err := base64.StdEncoding.DecodeString(strings.TrimSpace(data_str))
			if err != nil {
				log.Println("获取上传数据失败", err)
				WSErrHandler("发送文件失败，请重试", conn)
				continue
			}
			//reader := bytes.NewReader(temp)

			fi := msg.File
			//第一片的时候在服务器创建文件
			if 0 == fi.Segment {
				//获取保存到服务器的文件key
				key = model.GetFileKey(fi.Name)
				if key == 0 {
					WSErrHandler("服务器错误，请重试", conn)
					continue
				}
				log.Println("得到key为", key, "准备创建文件")
				cur_file_name = "../upload/" + strconv.FormatInt(key, 10)
				file, err = os.Create(cur_file_name)
				if err != nil {
					log.Println("创建文件失败", err)
					WSErrHandler("服务器错误，请重试", conn)
					continue
				}
			} else {
				file, err = os.OpenFile(cur_file_name, os.O_RDWR|os.O_APPEND, 0666)
				if err != nil {
					log.Println("打开文件失败", err)
					WSErrHandler("服务器错误，请重试", conn)
					continue
				}
			}

			if fi.Segment == fi.Total {
				//读取结束
				log.Println("上传文件到服务器成功")
				file.Close()

				msg.Value = "/upload/" + cur_file_name
			} else {
				//读取未结束
				_, err = file.Write(temp)
				//_, err = io.Copy(file, reader)
				if err != nil {
					log.Println("Copy", err)
					WSErrHandler("服务器错误，请重试", conn)
					continue
				}
				//读取未结束
				continue
			}

		case model.VIDEO:
			//bug 对方未登陆处理
			log.Println("接收方id:", msg.To)
			if !model.IsLogin(msg.To) {
				log.Println(msg.To, "未登陆,无法使用视频功能")
				//拒绝消息
				msg.Type = model.REJECT
				msg.Value = "用户未登陆，无法使用视频功能"
				temp := msg.To
				msg.To = msg.From
				msg.From = temp
				goto sendMsg
			}

		default:
			log.Println("收到不支持的消息格式", msg)
			return
		}

		//如果接收方没有连接信息，写入消息队列
		if Conns[msg.To] == nil {
			//此处应写入消息队列
			log.Println("由于用户", msg.To, "还未登陆，写入消息队列")
			model.AddMessage(msg.To, msg)
			continue
		}
	sendMsg:
		json_msg, err := json.Marshal(msg)
		if err != nil {
			log.Println("打包json失败")
			return
		}
		//转发
		Conns[msg.To].WriteMessage(websocket.TextMessage, json_msg)
		log.Println("转发消息成功")
	}
}
