package oss
import  (
  cfg "filestore_server/config"
  "fmt"
  "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

// 获取OSS客户端
func Client() *oss.Client {
  if ossCli == nil {
    ossCli, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccessKeyID, cfg.OSSAccessKeySecret)
	if err != nil {
	  fmt.Println("Error:", err)
	  return nil
	}
	return ossCli
  }
  return ossCli
}

// 获取bucket存储空间
func Bucket() *oss.Bucket{
	cli := Client()
	if(cli != nil){
		bucket, err:=cli.Bucket(cfg.OSSBucket)
		if(err != nil){
			fmt.Println("Error:", err)
			return nil
		}else{
			return bucket
		}
	}
	return nil
}

// 临时授权下载
func DownloadURL(objName string)string{
	signedUrl, err := Bucket().SignURL(objName, oss.HTTPGet, 60*60)
	if(err != nil){
		fmt.Println(err.Error())
		return ""
	}
	return signedUrl
}