package mysql

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init(){
	db,_ = sql.Open("mysql","root:123@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	// 同时活跃的连接数
	db.SetMaxOpenConns(100)
	err := db.Ping()
	if(err != nil){
		fmt.Println("Failed to connect to database, err: "+err.Error())
		os.Exit(1)
	}
}
// DBConn: 返回数据库连接对象
func DBConn() *sql.DB{
	return db
}