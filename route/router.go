package route

import (
	"filestore_server/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// gin framework, 包括Logger，Recovery
	router := gin.Default()
	// 处理静态资源
	router.Static("/static", "./static")
	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	router.GET("/user/signup", handler.SigninHandler)
	router.POST("/user/signup", handler.DoSigninHandler)

	router.POST("/user/info", handler.UserinfoHandler)
	router.POST("/file/meta",handler.GetFileMetaHandler)
	// 下载接口返回的data可能有问题
	router.POST("/file/download",handler.DownloadHandler)
	

	// 加入中间件, 用于校验token的拦截器


	// 这行代码之后的所有handler都会被这个中间件拦截


}
