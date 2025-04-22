package handler

import (
	"fmt"
	"io"
	"net/http"

	// 1.16之后，ioutil被弃用，改为os
	"encoding/json"
	"filestore_server/meta"
	"filestore_server/util"
	"os"
	"time"
)

// 处理文件上传服务
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Handle file upload
	if(r.Method == "GET"){
		// 返回上传的html界面
		data,err := os.ReadFile("./static/view/index.html")
		if(err != nil){
			io.WriteString(w,"Internal Server Error")
			return
		}
		io.WriteString(w,string(data))
	}else if(r.Method == "POST"){
		// 接收文件流存储到本地目录
		// 返回文件流（multipart.File）、文件头信息（*multipart.FileHeader）和错误
		file,head,err := r.FormFile("file")
		if(err != nil){
			fmt.Printf("Failed to get form file, err:%s\n", err.Error())
			return
		}
		// 直到 return 前才被执行。因此，可以用来做资源清理
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/"+head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if(err != nil){
			fmt.Printf("Failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if(err != nil){
			fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		// 如果程序无动作， go 处理完 http 会自动返回 200
		// http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}
}

// 上传成功的信息
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload success!")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	// 输入网址: http://124.223.141.236:8080/file/meta?filehash=04bad22b7045d31e34d1ca46cb03bbfcfd9d9fdc
	filehash := r.Form["filehash"][0]
	// fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	// 文件很小, 一次性读
	// TODO: 大文件分块读
	data,err := io.ReadAll(f)
	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// 文件重命名
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	// fmt.Printf("method: %s, op: %s, filehash: %s, filename: %s\n", r.Method, opType,fileSha1,newFileName)

	if(opType != "0"){
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if(r.Method != "POST"){
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data,err := json.Marshal(curFileMeta)

	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件及其元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}