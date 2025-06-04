package handler

import (
	dblayer "filestore_server/db"
	"filestore_server/util"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const(
	pwd_salt = "*#890"
)

// 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method == http.MethodGet){
		data, err := os.ReadFile("./static/view/signup.html")
		if(err != nil){
			log.Println("读取signup.html文件失败")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if(len(username) < 3 || len(passwd) < 5){
		w.Write([]byte("Invalid parameter, username or password too short"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if(suc){
		w.Write([]byte("SUCCESS"))
	}else{
		w.Write([]byte("FAILED"))
		return
	}
}

// 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request){
	log.Println("读取signin.html文件成功")
	if(r.Method == http.MethodGet){
		data, err := os.ReadFile("./static/view/signin.html")
		if(err != nil){
			log.Println("读取signin.html文件失败")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	// 1.校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, enc_passwd)
	if(!pwdChecked){
		w.Write([]byte("FAILED"))
		return
	}
	// 2.生成访问凭证token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if(!upRes){
		w.Write([]byte("FAILED"))
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg: "SUCCESS",
		Data: struct{
			Location string
			Username string
			Token string
		}{
			Location: "http://"+r.Host+"/static/view/home.html",
			Username: username,
			Token: token,
		},
	}
	w.Write(resp.JSONBytes())

	// 3.登录成功后重定向到首页
	// fmt.Println("user signed in succeed, web will jump")
	// w.Write([]byte("http://"+r.Host+"/static/view/home.html"))
}
// 查询用户信息
func UserinfoHandler(w http.ResponseWriter, r *http.Request){
	// 1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	// 2.验证token是否有效
	isValid := isTokenValid(token)
	if(!isValid){
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 3.获取用户信息
	user, err := dblayer.GetUserInfo(username)
	if(err != nil){
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4.组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: user,
	}

	w.Write(resp.JSONBytes())
}


func isTokenValid(token string) bool {
	// TODO:判断token的失效期，判断是否过期
	// TODO:从数据库中查询username对应的token信息
	// TODO:判断token是否匹配
	if(len(token) != 40){
		return false
	}
	return true
}

func GenToken(username string) string{
	// 40位字符: md5(username+timestamp+token_salt) + timestamp[0:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix+ts[:8]
}