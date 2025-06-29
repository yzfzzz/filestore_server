package main

import (
	"filestore_server/route"
	"log"
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
	r := route.Router()
	r.Run(":8080")

	// 高级功能
	// http.HandleFunc("/file/fastupload", handler.TryFastUploadHandler)
	// TODO: 分块上传这里有点问题
	// http.HandleFunc("/file/mpupload", handler.InitialMutipartUploadHandler)
	// http.HandleFunc("/file/mpupload/uppart", handler.UploadPartHandler)
	// http.HandleFunc("/file/mpupload/complete", handler.CompleteUploadHandler)
}