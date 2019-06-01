package dbops

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var (
	dbConn *sql.DB
	err error
)

func init(){
	connStr := "user=postgres dbname=video_server sslmode=disable"
	dbConn, err = sql.Open("postgres", connStr)
	//dbConn,err = sql.Open("mysql","root:root@tcp(localhost:3306)/video_server?charset=utf8")
	if err != nil{
		log.Fatalf("创建链接失败:%s",err)
	}
}
