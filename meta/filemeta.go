package meta
import(
	mydb "filestore_server/db"
)

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
// 将文件的元信息加入到mysql数据库中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 利用sha1获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}
// 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error){
	tfile,err := mydb.GetFileMeta(fileSha1)
	if(err != nil){
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

// 删除元信息
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}

func IsEmpty(f *FileMeta) bool {
    return f.FileName == "" && f.FileSize == 0 && f.FileSha1  == ""
}
