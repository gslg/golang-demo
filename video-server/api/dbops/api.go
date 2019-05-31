package dbops

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)


func AddUserCridential(loginName string,pwd string) error {
	sql := "INSERT INTO users (`login_name`,`password`) values (?,?)"

	stmt, e := dbConn.Prepare(sql)

	if e != nil{
		log.Printf("添加用户报错:%s",e)
		return e
	}
	defer stmt.Close()
	result, e := stmt.Exec(loginName, pwd)
	log.Println(result.LastInsertId())
	return e
}

func GetUserCridential(loginName string) (string,error) {
	stmt,err := dbConn.Prepare("SELECT `password` from user where login_name = ?")

	if err != nil{
		log.Printf("查询用户报错:%s",err)
		return "",err
	}
	defer stmt.Close()
	var pwd string
	stmt.QueryRow(loginName).Scan(&pwd)

	return pwd,nil
}

func DeleteUser(loginName string,pwd string) error {
	stmt,err := dbConn.Prepare("DELETE FROM user where login_name = ? and password = ? ")
	if err != nil{
		log.Printf("删除用户报错:%s",err)
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(loginName, pwd)

	log.Println(result)
	return err
}

