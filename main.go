package main

import (
	"filestore_server/handler"
	"fmt"
	"net/http"
)
func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	err := http.ListenAndServe(":8080", nil)
	if(err != nil){
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}