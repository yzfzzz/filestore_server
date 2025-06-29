package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	// 1.16之后，ioutil被弃用，改为os
	"encoding/json"
	cmn "filestore_server/common"
	cfg "filestore_server/config"
	dblayer "filestore_server/db"
	"filestore_server/meta"
	mq "filestore_server/mq"
	oss "filestore_server/store"
	"filestore_server/util"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// 处理文件上传服务
func UploadHandler(c *gin.Context) {
	// Handle file upload
	c.Redirect(http.StatusFound, "/static/view/index.html")
}

func DoUploadHandler(c *gin.Context) {
	// 接收文件流存储到本地目录
	// 返回文件流（multipart.File）、文件头信息（*multipart.FileHeader）和错误
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get form file, err:%s\n", err.Error())
		return
	}
	// 直到 return 前才被执行。因此，可以用来做资源清理
	defer file.Close()
	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "./tmp/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err:%s\n", err.Error())
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	// 写入阿里云oss存储
	newFile.Seek(0, 0)

	ossPath := "oss/" + fileMeta.FileSha1
	// 判断同步还是异步
	if !cfg.AsyncTransferEnable {
		err = oss.Bucket().PutObject(ossPath, newFile)
		if err != nil {
			fmt.Printf("Failed to upload data into OSS, err:%s\n", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"msg":  "Failed to upload data into OSS",
				"code": -1,
			})
			return
		}
		fileMeta.Location = ossPath
	} else {

		data := mq.TransferData{
			FileHash:      fileMeta.FileSha1,
			CurLocation:   fileMeta.Location,
			DestLocation:  ossPath,
			DestStoreType: cmn.StoreOSS,
		}

		pubData, _ := json.Marshal(data)
		pubSuc := mq.Publish(
			cfg.TransExchangeName,
			cfg.TransOSSRoutingKey,
			pubData,
		)
		if !pubSuc {
			// 发送消息失败
			log.Println("send msg error!")
		}
	}

	// meta.UpdateFileMeta(fileMeta)
	_ = meta.UpdateFileMetaDB(fileMeta)

	// 更新用户文件表记录
	username := c.Request.FormValue("username")

	suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
	if suc == true {
		c.Redirect(http.StatusOK, "/static/view/home.html")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Upload failed!",
			"code": -1,
		})
	}

	// 如果程序无动作， go 处理完 http 会自动返回 200
	// http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

}

// 上传成功的信息
func UploadSucHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  "Upload success!",
	})
}

// 获取文件元信息
func GetFileMetaHandler(c *gin.Context) {
	// 输入网址: http://124.223.141.236:8080/file/meta?filehash=04bad22b7045d31e34d1ca46cb03bbfcfd9d9fdc
	filehash := c.Request.FormValue("filehash")

	// fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to get file meta data",
			"code": -1,
		})
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to json",
			"code": -2,
		})
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")

	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to query user file metas",
			"code": -1,
		})
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Failed to json",
			"code": -2,
		})
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

// 文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to open file",
			"code": -1,
		})
		return
	}
	defer f.Close()
	// 文件很小, 一次性读
	// TODO: 大文件分块读
	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to read file",
			"code": -2,
		})
		return
	}

	c.Writer.Header().Add("Content-Type", "application/octect-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+fm.FileName+"\"")
	c.Data(http.StatusOK, "application/octect-stream", data)
}

// 文件重命名
func FileMetaUpdateHandler(c *gin.Context) {
	

	opType := c.Request.FormValue("op")
	fileSha1 := c.Request.FormValue("filehash")
	newFileName := c.Request.FormValue("filename")

	// fmt.Printf("method: %s, op: %s, filehash: %s, filename: %s\n", r.Method, opType,fileSha1,newFileName)

	if opType != "0" {
		c.JSON(http.StatusForbidden, gin.H{
			"msg":  "Forbidden",
			"code": -1,
		})
		return
	}


	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to json",
			"code": -2,
		})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

// 删除文件及其元信息
func FileDeleteHandler(c *gin.Context) {
	fileSha1 := c.Request.FormValue("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "Delete success!",
	})
}

// 尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 1.解析请求参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 2.从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Printf("Failed to query file hash, err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	// 3.查无记录，返回秒传失败
	if meta.IsEmpty(&fileMeta) {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// 4.有记录，则将文件信息写入用户文件表，返回成功
	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
	} else {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
		}
		w.Write(resp.JSONBytes())
		return
	}
}

func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	// 从文件表查找记录
	row, err := dblayer.GetFileMeta(filehash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Failed to get file meta data",
			"code": -1,
		})
		return
	}
	//
	if row.FileAddr.Valid {
		signedURL := oss.DownloadURL(row.FileAddr.String)
		c.Data(http.StatusOK, "text/plain", []byte(signedURL))
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to get file download url",
			"code": -2,
		})
	}
}
