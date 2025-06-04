package main

import (
	"bufio"
	"encoding/json"
	"filestore_server/config"
	dblayer "filestore_server/db"
	"filestore_server/mq"
	oss "filestore_server/store"
	"log"
	"os"
)

func ProcessTransfer(msg []byte)bool{
	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil{
		log.Println(err.Error())
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil{
		log.Println(err.Error())
		return false
	}

	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin))
	
	if err != nil{
		log.Println(err.Error())
		return false
	}

	_ = dblayer.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)

	return true
}

func main(){
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if !config.AsyncTransferEnable{
		log.Println("异步转移文件功能未开启")
		return
	}
	log.Println("文件转移服务启动, 开始监听")
	mq.StartConsumer(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}