package controller

import (
	"log"
	"model"
	"net/http"
	"strconv"
	"tool"

	"github.com/gorilla/mux"
)

//用户相关逻辑处理，
//Get：解析页面
//Post：注册用户
//Put：修改用户信息
func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: //Get请求，解析register.html页面
		UserGetHandler(w, r)
	case http.MethodPost: //用户注册处理
	case http.MethodDelete:
	case http.MethodPut:
	default: //不支持别的方法
	}
}

//解析register.html页面
func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if "" == vars["uid"] { //为空说明为"/user"，解析register页面
		if err := RenderTmpl(TMPL_REGISTER, nil, w); err != nil {
			//500错误
		}
	} else { //否则为获取用户信息

	}
}

//注册逻辑
func UserPostHanlder(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("uname")
	password := r.FormValue("password")
	repassword := r.FormValue("repassword")
	motto := r.FormValue("motto")
	email := r.FormValue("email") //目前的机制由于找回密码需要必须填写邮箱
	portrait := "default.jpg"     //暂时写死

	//检查数据合法性
	if !tool.ValidateUname(uname) || !tool.ValidatePsw(password) || !tool.ValidateEmail(email) {
		w.Write([]byte("用户名或密码格式不合法"))
		return
	}
	if password != repassword {
		w.Write([]byte("两次输入的密码不相同"))
		return
	}
	//bug邮箱校验

	user := model.NewUser(uname, password, motto, portrait, email)
	uid := user.CreateUser()
	if 0 == uid {
		w.Write([]byte("注册失败，请重新注册"))
		return
	}
	log.Println("用户", uid, "注册成功")

	http.Redirect(w, r, "/success?ref=register&uid="+strconv.FormatInt(uid, 10), http.StatusFound)
}

//更新用户信息
func UserPutHanlder(w http.ResponseWriter, r *http.Request) {

}
