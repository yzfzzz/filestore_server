package main

import (
	"filestore_server/handler"
	"fmt"
	"net/http"
)
func main() {
	// 静态资源处理
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.AssetFS())))
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	// TODO: 在路由处加入拦截器
	http.HandleFunc("/user/info", handler.UserinfoHandler)
	err := http.ListenAndServe(":8080", nil)
	if(err != nil){
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}