package handler

import (
	dblayer "filestore_server/db"
	"filestore_server/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	pwd_salt = "*#890"
)

// 处理用户注册get请求
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// 处理注册post请求
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Invalid parameter, username or password too short",
			"code": -1,
		})
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Signup succeed",
			"code": 0,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Signup failed",
			"code": -2,
		})
		return
	}
}

// 登录接口
func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// 处理登录Post接口
func DoSigninHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	// 1.校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, enc_passwd)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Password incorrect",
			"code": -1,
		})
		return
	}
	// 2.生成访问凭证token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Token incorrect",
			"code": -2,
		})
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

// 查询用户信息
func UserinfoHandler(c *gin.Context) {
	// 1.解析请求参数
	username := c.Request.FormValue("username")
	token := c.Request.FormValue("token")
	// 2.验证token是否有效
	isValid := isTokenValid(token)
	if !isValid {
		c.JSON(http.StatusForbidden, gin.H{
			"msg":  "Token is not valid, Forbidden acess!",
			"code": -1,
		})
		return
	}

	// 3.获取用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"msg":  "Usename is not valid, Forbidden acess!",
			"code": -2,
		})
		return
	}

	// 4.组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}

	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

func isTokenValid(token string) bool {
	// TODO:判断token的失效期，判断是否过期
	// TODO:从数据库中查询username对应的token信息
	// TODO:判断token是否匹配
	if len(token) != 40 {
		return false
	}
	return true
}

func GenToken(username string) string {
	// 40位字符: md5(username+timestamp+token_salt) + timestamp[0:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
