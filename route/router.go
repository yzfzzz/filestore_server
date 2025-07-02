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
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.DoSigninHandler)

	router.POST("/user/info", handler.UserinfoHandler)
	router.POST("/file/meta",handler.GetFileMetaHandler)
	// 下载接口返回的data可能有问题
	router.POST("/file/download",handler.DownloadHandler)
	router.GET("/file/upload",handler.UploadHandler)
	router.POST("/file/upload",handler.DoUploadHandler)
	router.POST("/file/upload/suc",handler.UploadSucHandler)
	router.POST("/file/update",handler.FileMetaUpdateHandler)
	router.POST("/file/delete", handler.FileDeleteHandler)
	router.POST("/file/query", handler.FileQueryHandler)
	router.POST("/file/downloadurl", handler.DownloadURLHandler)

	// 加入中间件, 用于校验token的拦截器


	// 这行代码之后的所有handler都会被这个中间件拦截

	return router
}
