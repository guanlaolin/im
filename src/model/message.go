package model

import (
	"db"
	"encoding/json"
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
)

//支持的消息格式
const (
	SUCCESS  = iota //成功0
	ERROR           //错误1
	TEXT            //文本2
	FILE            //文件3
	VIDEO           //视频4
	IDENTITY        //首次连接确认身份5
	REJECT          //请求拒绝6
)

//消息结构体
type Message struct {
	Type  int         `json:"type"`
	From  int         `json:"from"`
	To    int         `json:"to"`
	Value interface{} `json:"value"`
	File  FileInfo    `json:"fileinfo"`
	Back1 interface{} `json:"back1"`
}

//未读消息结构体
type UnReadMsg struct {
	Num   int       `json:"num"`   //消息数量
	Value []Message `json:"value"` //消息
}

//上传文件信息
type FileInfo struct {
	Name    string `json:"name"`    //文件名
	Segment int    `json:"segment"` //当前是第几片
	Total   int    `json:"total"`   //总片数
}

//写入消息队列
func AddMessage(uid int, msg Message) {
	conn := db.GetNoConn()
	if conn == nil {
		return
	}
	defer conn.Close()
	v, err := json.Marshal(msg)
	if err != nil {
		log.Println("Marshal", err)
		return
	}
	_, err = conn.Do("LPUSH", uid, string(v))
	if err != nil {
		log.Fatalln(err)
	}
}

//获取消息队列中所有信息并清空消息队列
func GetAllMessage(uid int) []Message {
	//连接redis
	conn := db.GetNoConn()
	if conn == nil {
		return nil
	}
	defer conn.Close()

	//获取消息队列中所有消息
	temp, err := redis.Values(conn.Do("LRANGE", uid, 0, -1))
	if err != nil {
		log.Println("LRANGE", err)
		return nil
	}

	var msgs []Message
	for _, v := range temp {
		//去除redis存储添加的'\'
		rpStr := strings.Replace(string(v.([]byte)), "\\", "", -1)
		var msg Message
		err := json.Unmarshal([]byte(rpStr), &msg)
		if err != nil {
			log.Println("Unmarshal", err)
			return nil
		}
		msgs = append(msgs, msg)
	}

	//清空消息队列
	_, err = conn.Do("DEL", uid)
	if err != nil {
		log.Println("LTRIM", err)
		return nil
	}

	return msgs
}
