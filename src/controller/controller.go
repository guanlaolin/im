package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"logger"
	"model"
	"net/http"
	"strconv"
	"time"
	"tool"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

//模板对应文件名
const (
	TMPL_DIR      = "../tmpl/"
	TMPL_INDEX    = TMPL_DIR + "index.html"
	TMPL_REGISTER = TMPL_DIR + "register.html"
	TMPL_LOGIN    = TMPL_DIR + "login.html"
)

//日志初始化
var l = logger.NewLoggerWithFile(logger.LEVEL_DEBUG, "../log/controller_log")

//session
var store = sessions.NewCookieStore([]byte(tool.Conf("session-secret")))

//保存websocket连接
var Conns = make(map[int]*websocket.Conn)

//记录uuid
var UUIDS []byte

//主聊天页面逻辑
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		IndexGetHandler(w, r)
	default:
		//不支持的方法
	}
}

//解析index.html页面
func IndexGetHandler(w http.ResponseWriter, r *http.Request) {
	//	datas := make(map[string]interface{})

	//判断是否已经登录
	session, err := store.Get(r, tool.Conf("session-name"))
	if err != nil {
		//500 错误
	}

	if session.IsNew {
		//未登陆，跳转到登录页面，还需优化
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	/*
		uid := session.Values["uid"]
		if uid != "" {
			//未登陆，跳转到登录页面
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	*/

	//用户好友信息考虑调用接口，这里只解析index.html
	if err = RenderTmpl(TMPL_INDEX, nil, w); err != nil {
		//500
	}

	//	loginUser := model.NewLoginUser(uid, "")

	//	//获取用户本身概要信息
	//	smu := loginUser.GetUser(uid)
	//	log.Println("用户信息", smu)
	//	datas["ownInfo"] = smu

	//	//获取好友列表
	//	friends := loginUser.GetFriends()
	//	log.Println("好友列表", friends)
	//	datas["friendList"] = friends

	//	tmpl, err := template.ParseFiles(TMPL_DIR + "/index.html")
	//	if err != nil {
	//		log.Println("解析页面index.html失败", err.Error())
	//		http.Error(w, "服务器错误，请刷新", http.StatusInternalServerError)
	//		return
	//	}
	//	if err = tmpl.Execute(w, datas); err != nil {
	//		log.Println("渲染模板失败", err.Error())
	//	}
	//	log.Println("渲染模板index.html成功")
}

//查找用户
func SearchUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("uid")
		if err != nil {
			log.Println("读取cookie失败，用户未登陆", err)
			//说明未登陆，跳转
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		uid, _ := strconv.Atoi(cookie.Value)
		loginUser := model.NewLoginUser(uid, "")

		fid, err := strconv.Atoi(r.FormValue("fid"))
		if err != nil || !model.ValidateUid(fid) {
			log.Println("上送的uid不合法", err)
			w.Write([]byte("请输入合法的用户id"))
			return
		}

		smUser := loginUser.GetUser(fid)
		log.Println("查找到用户:", smUser)
		//bug
		if 0 == smUser.Uid {
			w.Write([]byte("empty"))
			return
		}
		//此处基本不可能出现错误，不处理error
		json_data, _ := json.Marshal(smUser)
		w.Write(json_data)

	} else {
		http.Error(w, "未实现的方法", http.StatusMethodNotAllowed)
		return
	}
}

//获取用户信息
func UserInfoHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("uid")
		if err != nil {
			log.Println("读取cookie失败，用户未登陆", err)
			//说明未登陆，跳转
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		uid, _ := strconv.Atoi(cookie.Value)
		loginUser := model.NewLoginUser(uid, "")

		fid, err := strconv.Atoi(r.FormValue("fid"))
		if err != nil || !model.ValidateUid(fid) {
			log.Println("上送的uid不合法", err)
			w.Write([]byte("请输入合法的用户id"))
			return
		}

		user := loginUser.GetUserInfo(fid)
		log.Println("查找到用户:", user)
		//bug
		if 0 == user.Uid {
			w.Write([]byte("empty"))
			return
		}
		//此处基本不可能出现错误，不处理error
		json_data, _ := json.Marshal(user)
		w.Write(json_data)

	} else {
		http.Error(w, "未实现的方法", http.StatusMethodNotAllowed)
		return
	}
}

//添加好友逻辑
func AddFriendHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("uid")
		if err != nil {
			log.Println("读取cookie失败，用户未登陆", err)
			//说明未登陆，跳转
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		uid, _ := strconv.Atoi(cookie.Value)
		loginUser := model.NewLoginUser(uid, "")

		fid, _ := strconv.Atoi(r.FormValue("fid"))
		log.Println(fid)
		if !model.ValidateUid(fid) {
			w.Write([]byte("您输入的id不合法"))
			return
		}

		if !loginUser.AddFriend(fid) {
			log.Println("添加好友失败")
			w.Write([]byte("添加好友失败，请重试"))
			return
		}
		w.Write([]byte("success"))
	} else {
		http.Error(w, "未实现的方法", http.StatusMethodNotAllowed)
		return
	}
}

//删除好友逻辑
func DelFriendHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("uid")
		if err != nil {
			log.Println("读取cookie失败，用户未登陆", err)
			//说明未登陆，跳转
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		uid, _ := strconv.Atoi(cookie.Value)
		loginUser := model.NewLoginUser(uid, "")

		fid, _ := strconv.Atoi(r.FormValue("fid"))
		log.Println(fid)
		if !model.ValidateUid(fid) {
			w.Write([]byte("您输入的id不合法"))
			return
		}

		if !loginUser.DelFriend(fid) {
			log.Println("删除好友失败")
			w.Write([]byte("删除好友失败，请重试"))
			return
		}
		log.Println(uid, "删除好友", fid, "成功")
		w.Write([]byte("success"))
	} else {
		http.Error(w, "未实现的方法", http.StatusMethodNotAllowed)
		return
	}
}

//未读消息逻辑
func UnReadMsgHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var unReadMsg model.UnReadMsg
		var msgs []model.Message

		unReadMsg.Num = -1

		uid, err := strconv.Atoi(r.FormValue("uid"))
		if err != nil {
			log.Println("上送的uid有误，获取未读消息失败")
			goto send_msg
		}
		if !model.IsLogin(uid) {
			log.Println("用户未登陆")
			goto send_msg
		}

		msgs = model.GetAllMessage(uid)
		unReadMsg.Num = len(msgs)
		unReadMsg.Value = msgs

	send_msg:
		json_data, err := json.Marshal(unReadMsg)
		if err != nil {
			log.Println("打包json失败")
		}
		_, err = w.Write(json_data)
		if err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "未实现的方法", http.StatusMethodNotAllowed)
		return
	}
}

//重置密码
//bug:未做校验
func UpdatePswHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var user *model.LoginUser

		uid, err := strconv.Atoi(r.FormValue("uid"))
		if err != nil || !model.ValidateUid(uid) {
			w.Write([]byte("请输入合法的uid"))
			return
		}
		new_password := r.FormValue("newpassword")
		re_new_password := r.FormValue("newrepassword")
		if !model.ValidatePsw(new_password) {
			w.Write([]byte("请输入合法的密码"))
			return
		}
		if new_password != re_new_password {
			log.Println("两次密码不一致")
			w.Write([]byte("两次输入的密码不一样"))
			return
		}

		//判断是修改密码还是重置密码，或者后续添加更多功能
		switch r.FormValue("type") {
		case "change": //修改密码
			//修改密码需要校验旧密码
			old_password := r.FormValue("oldpassword")
			if !model.ValidatePsw(old_password) {
				w.Write([]byte("请输入合法的密码"))
				return
			}
			user = model.NewLoginUser(uid, old_password)
			if !user.ValidateUserPassword() {
				w.Write([]byte("原账号密码不匹配，修改密码失败"))
				return
			}
		case "reset": //重置密码
			//重置密码需要校验邮箱
			email := r.FormValue("email")
			if !model.ValidateEmail(email) {
				w.Write([]byte("邮箱不合法"))
				return
			}
			user = model.NewLoginUser(uid, "")
			if !user.ValidateUserEmail(email) {
				w.Write([]byte("校验邮箱失败，重置密码失败"))
				return
			}
		default:
			log.Println("不支持的操作")
			w.Write([]byte("不支持的修改密码操作"))
			return
		}
		if !user.UpdatePassword(new_password) {
			w.Write([]byte("服务器错误，修改密码失败，请稍后重试"))
			return
		}
		w.Write([]byte("success"))
		log.Println(uid, "修改密码成功")
	} else {
		http.Error(w, "未实现的方法", http.StatusMethodNotAllowed)
		return
	}
}

//成功
func SuccessHandle(w http.ResponseWriter, r *http.Request) {
	var renderData = make(map[string]interface{})
	tmpl, err := template.ParseFiles(TMPL_DIR + "/success.html")
	if err != nil {
		log.Println("解析页面success.html失败", err.Error())
		http.Error(w, "服务器错误，请重试", http.StatusInternalServerError)
		return
	}
	switch r.FormValue("ref") {
	case "register":
		renderData["title"] = "注册成功"
		renderData["message"] = template.HTML(`<h3>恭喜你,注册成功，你的ID是：` + r.FormValue("uid") + ` 请牢记</h3>` +
			`<br /><a href="/login">点击此处跳转到登录页面</a>`)
	default:
		renderData["message"] = `<h3>引用源有误</h3>`
	}
	//渲染数据
	err = tmpl.Execute(w, renderData)
	if err != nil {
		log.Println("渲染register.html失败", err)
		return
	}
}

//这里简单实现，检测websocket客户端是否异常断开，应完善机制
func Ping() {
	for {
		time.Sleep(5 * time.Second)
		for k, _ := range Conns {
			err := Conns[k].WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				delete(Conns, k)
				l.DEBUG(k, "已经从服务器断开")
				continue
			}
			l.DEBUG("往", k, "发送心跳包成功")
		}
	}
}
