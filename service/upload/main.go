package main

import (
	"filestore_server/handler"
	"fmt"
	"log"
	"net/http"
	"os"
)

func logInit() *os.File{
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	file, err := os.OpenFile("/home/filestore_server/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法打开日志文件:", err)
	}
	log.SetOutput(file)
	return file
}

func main() {
	file := logInit()
	defer file.Close()
	// 静态资源处理
	http.Handle("/static/",http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// TODO: 在路由处加入拦截器
	// 用户模块
	http.HandleFunc("/", handler.SignInHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.UserinfoHandler)

	// 文件管理
	// 基本功能
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/downloadurl", handler.DownloadURLHandler)

	// 高级功能
	http.HandleFunc("/file/fastupload", handler.TryFastUploadHandler)

	// TODO: 分块上传这里有点问题
	// http.HandleFunc("/file/mpupload", handler.InitialMutipartUploadHandler)
	// http.HandleFunc("/file/mpupload/uppart", handler.UploadPartHandler)
	// http.HandleFunc("/file/mpupload/complete", handler.CompleteUploadHandler)

	err := http.ListenAndServe(":8080", nil)
	if(err != nil){
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}