package model

import (
	"database/sql"
	"db"
	"log"
)

type User struct {
	Uid      int    `json:"uid"`
	Uname    string `json:"uname"` //昵称
	Password string `json:"password"`
	Motto    string `json:"motto"`
	Portrait string `json:"portrait"` //头像地址
	Email    string `json:"email"`
}

type SummaryUser struct {
	Uid      int    `json:"uid"`      //uid
	Uname    string `json:"uname"`    //昵称
	Motto    string `json:"motto"`    //个性签名
	Portrait string `json:"portrait"` //头像地址
	Login    bool   `json:"login"`    //是否已经登录
}

type LoginUser struct {
	Uid      int
	Password string
}

func NewUser(uname string, psw string, motto string, portrait string, email string) *User {
	return &User{Uname: uname, Password: psw, Motto: motto, Portrait: portrait, Email: email}
}

func NewLoginUser(uid int, psw string) *LoginUser {
	return &LoginUser{uid, psw}
}

//创建用户（用户注册),返回创建的用户ID
//未清洗数据
func (u *User) CreateUser() int64 {
	//检查数据合法性
	if !u.checkUser() {
		return 0
	}

	//清洗数据

	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return 0
	}
	defer d.Close()

	//插入数据库
	result, err := d.Exec("INSERT INTO tb_user(uname,password,motto,portrait,email) VALUES(?,?,?,?,?)",
		u.Uname, encrypt(u.Password), u.Motto, u.Portrait, u.Email)
	if err != nil {
		log.Println("插入数据库失败")
		return 0
	}
	uid, _ := result.LastInsertId()
	return uid
}

//用户登录
func (u *LoginUser) Login() bool {
	if !u.checkLoginUser() {
		return false
	}

	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return false
	}
	defer d.Close()

	//查询数据库
	row := d.QueryRow("SELECT password FROM tb_user WHERE uid = ?", u.Uid)

	var psw string
	row.Scan(&psw)
	if psw != encrypt(u.Password) {
		log.Println("用户名或密码错误")
		return false
	}
	//写入redis
	conn := db.GetNoConn()
	if conn == nil {
		return false
	}

	defer conn.Close()

	_, err := conn.Do("SADD", "LOGINED", u.Uid)
	if err != nil {
		log.Println("插入redis失败", err.Error())
		return false
	}

	return true
}

//用户注销
func (u *LoginUser) Logout() {
	if !u.IsLogin() {
		return
	}
	conn := db.GetNoConn()
	if conn == nil {
		return
	}
	conn.Do("SREM", "LOGINED", u.Uid)
}

//判断用户是否已经登录
func (u *LoginUser) IsLogin() bool {
	return IsLogin(u.Uid)
}

//添加好友
func (u *LoginUser) AddFriend(uid int) bool {
	if !u.checkLoginUser() || !CheckUid(uid) {
		return false
	}

	//判断是否已经登录
	if !u.IsLogin() {
		return false
	}

	//判断是否已经是好友
	if !u.IsNotFriend(uid) {
		return false
	}

	if u.Uid == uid {
		log.Println("不能添加自己为好友")
		return false
	}

	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return false
	}
	defer d.Close()

	//插入数据库
	_, err := d.Exec("CALL proc_add_friend(?,?)", u.Uid, uid)
	if err != nil {
		log.Println("插入数据库失败，添加好友失败", err.Error())
		return false
	}

	return true
}

//判断是否已经是好友，是返回false，不是返回true，错误返回false
func (u *LoginUser) IsNotFriend(fid int) bool {
	if !u.checkLoginUser() || !CheckUid(fid) {
		return false
	}

	//判断是否已经登录
	if !u.IsLogin() {
		return false
	}

	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return false
	}
	defer d.Close()

	//查找数据库判断是否已经是好友
	//bug
	var temp string
	err := d.QueryRow("SELECT fid FROM tb_friend WHERE uid=? AND fid=?", u.Uid, fid).Scan(&temp)
	if err != nil {
		//log.Println("插入数据库失败，添加好友失败", err.Error())
		log.Println("应该还不是好友", err)
		return true
	}

	return false
}

//删除好友
func (u *LoginUser) DelFriend(uid int) bool {
	if !u.checkLoginUser() || !CheckUid(uid) {
		return false
	}

	//判断是否已经登录
	if !u.IsLogin() {
		return false
	}

	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return false
	}
	defer d.Close()

	//插入数据库
	_, err := d.Exec("DELETE FROM tb_friend WHERE uid = ? AND fid = ?", u.Uid, uid)
	if err != nil {
		log.Println("删除好友失败，删除好友失败", err.Error())
		return false
	}

	return true
}

//用户是否已经登录
func IsLogin(uid int) bool {
	if !CheckUid(uid) {
		return false
	}

	conn := db.GetNoConn()
	if conn == nil {
		return false
	}
	defer conn.Close()

	res, err := conn.Do("SISMEMBER", "LOGINED", uid)
	if err != nil {
		log.Println("IsLogin:查询redis失败", err.Error())
		return false
	}

	if res == 0 {
		return false
	}

	return true
}

//获取好友列表
func (u *LoginUser) GetFriends() (friends []SummaryUser) {
	if !u.checkLoginUser() {
		return nil
	}

	//判断用户是否已经登录
	if !u.IsLogin() {
		return nil
	}

	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return nil
	}
	defer d.Close()

	//查询数据库
	rows, err := d.Query("SELECT uid,uname FROM tb_user WHERE uid = "+
		"ANY(SELECT fid FROM tb_friend WHERE uid = ?)", u.Uid)
	if err != nil {
		log.Println("查询数据库失败，获取好友列表失败", err.Error())
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var friend SummaryUser
		if err := rows.Scan(&friend.Uid, &friend.Uname); err != nil {
			return nil
		}

		if IsLogin(friend.Uid) {
			friend.Login = true
		} else {
			friend.Login = false
		}
		friends = append(friends, friend)
	}
	return friends
}

//获取单个用户简要信息
func (u *LoginUser) GetUser(uid int) (su SummaryUser) {
	if !u.checkLoginUser() || !CheckUid(uid) {
		return
	}
	if !u.IsLogin() {
		return
	}
	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return SummaryUser{}
	}
	defer d.Close()

	//查询数据库
	row := d.QueryRow("SELECT uid,uname,motto,portrait FROM tb_user WHERE uid = ?", uid)
	//防止空值问题
	var motto, portrait sql.NullString
	if err := row.Scan(&su.Uid, &su.Uname, &motto, &portrait); err != nil {
		log.Println("GetUser, scan error", err)
		return SummaryUser{}
	}
	if motto.Valid == true {
		su.Motto = motto.String
	}
	if portrait.Valid == true {
		su.Portrait = portrait.String
	}
	return su
}

//获取单个用户详细信息
func (u *LoginUser) GetUserInfo(uid int) (user User) {
	if !u.checkLoginUser() || !CheckUid(uid) {
		return
	}
	if !u.IsLogin() {
		return
	}
	//连接数据库
	d := db.GetDBConn()
	if d == nil {
		return User{}
	}
	defer d.Close()

	//查询数据库
	row := d.QueryRow("SELECT uid,uname,email,motto,portrait FROM tb_user WHERE uid = ?", uid)
	//防止空值问题
	var motto, portrait sql.NullString
	if err := row.Scan(&user.Uid, &user.Uname, &user.Email, &motto, &portrait); err != nil {
		log.Println("GetUserInfo, scan error", err)
		return User{}
	}
	if motto.Valid == true {
		user.Motto = motto.String
	}
	if portrait.Valid == true {
		user.Portrait = portrait.String
	}
	return user
}

//校验账号密码是否对应
//bug未做数据校验
func (u *LoginUser) ValidateUserPassword() bool {
	var password string

	d := db.GetDBConn()
	if nil == d {
		return false
	}
	err := d.QueryRow("SELECT password FROM tb_user WHERE uid=?", u.Uid).Scan(&password)
	if err != nil {
		log.Println("校验账号密码失败", err)
		return false
	}
	if encrypt(u.Password) != password {
		log.Println("账号密码不匹配")
		return false
	}
	return true
}

//校验账号邮箱是否对应
//bug未做数据校验
func (u *LoginUser) ValidateUserEmail(email string) bool {
	var mail string

	d := db.GetDBConn()
	if nil == d {
		return false
	}
	err := d.QueryRow("SELECT email FROM tb_user WHERE uid=?", u.Uid).Scan(&mail)
	if err != nil {
		log.Println("校验账号邮箱失败", err)
		return false
	}
	if email != mail {
		return false
	}
	return true
}

//修改密码
//未做数据校验
func (u *LoginUser) UpdatePassword(new_password string) bool {
	d := db.GetDBConn()
	if nil == d {
		return false
	}

	_, err := d.Exec("UPDATE tb_user SET password=? WHERE uid=?", encrypt(new_password), u.Uid)
	if err != nil {
		log.Println("修改密码失败", err)
		return false
	}
	return true
}
