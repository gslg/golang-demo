package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	dbConn *sql.DB
	err error
)

func init(){
	dbConn,err = sql.Open("mysql","root:root@tcp(localhost:3306)/video_server?charset=utf8")
	if err != nil{
		log.Fatalf("创建链接失败:%s",err)
	}
}
