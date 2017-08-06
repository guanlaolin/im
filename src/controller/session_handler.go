package controller

import (
	"model"
	"net/http"
	"strconv"
	"tool"
)

//回话管理
func SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: //Get请求，解析login.html页面
		SessionGetHandler(w, r)
	case http.MethodPost: //post请求，登录处理
		SessionPostHandler(w, r)
	case http.MethodDelete: //Delete请求，注销操作
		SessionDeleteHandler(w, r)
	default: //不支持别的方法
	}
}

//解析login.html页面
func SessionGetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, tool.Conf("session-name"))
	if err != nil {
		//500错误
	}

	//说明已经登录，自动跳转到主页面
	if session.Values["uid"] != "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := RenderTmpl(TMPL_LOGIN, nil, w); err != nil {
		//500错误
	}
}

//用户登录逻辑
func SessionPostHandler(w http.ResponseWriter, r *http.Request) {
	//获取uid
	uid, err := strconv.Atoi(r.FormValue("uid"))
	if err != nil {
		l.DEBUG("转换uid错误，uid：", r.FormValue("uid"), err)
		w.Write([]byte("请输入合法uid"))
		return
	}
	//获取用户密码
	password := r.FormValue("password")
	//判断是否记住用户
	remember := r.FormValue("remember")

	//检查数据合法性
	if !tool.ValidateUid(uid) || !tool.ValidatePsw(password) {
		w.Write([]byte("uid或密码格不合法"))
		return
	}

	loginUser := model.NewLoginUser(uid, password)
	if !loginUser.Login() {
		l.DEBUG("用户登录失败，uid:")
		w.Write([]byte("用户登录失败，请重试"))
		return
	}
	l.DEBUG("用户登录成功，uid:", uid)

	//session
	maxAge := 0
	if remember == "true" {
		//如果用户点了记住密码，设置cookie有效期为7天
		maxAge = 7 * 24 * 60 * 60
	}

	store.MaxAge(maxAge)
	session, err := store.Get(r, tool.Conf("session-name"))
	if err != nil {

	}
	session.Values["uid"] = uid
	if err := session.Save(r, w); err != nil {
		l.DEBUG("session save error", err)
		//登录失败
		w.Write([]byte("false"))
	}
}

//用户注销逻辑
func SessionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(r.FormValue("uid"))
	if err != nil {
		l.DEBUG("atoi,uid", uid, err)
		w.Write([]byte("上送的uid有误，注销失败，请重试"))
		return
	}

	model.NewLoginUser(uid, "").Logout()
	l.DEBUG("用户", uid, "注销成功")
	w.Write([]byte("success"))
}
