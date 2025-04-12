package meta

// 文件元信息结构
type FileMeta struct{
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	// 时间戳
	UploadAt string
}

// 类似于cpp的 map<string,FileMeta>
var fileMetas map[string]FileMeta

// 初始化
func init(){
	fileMetas = make(map[string]FileMeta)
}

// 更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1] = fmeta
}

// 利用sha1获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}

// 删除元信息
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}
