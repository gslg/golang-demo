package dbops

import (
	_ "github.com/lib/pq"
	"log"
	"time"
	"video-server/api/defs"
	"video-server/util"
)

func AddUserCredential(loginName string, pwd string) error {
	sql := `INSERT INTO "user" (login_name,password) values ($1,$2) RETURNING id`

	stmt, e := dbConn.Prepare(sql)

	if e != nil {
		log.Printf("添加用户报错: %s", e)
		return e
	}
	defer stmt.Close()
	id := 0
	e = stmt.QueryRow(loginName, pwd).Scan(&id)

	if e != nil {
		log.Printf("添加用户报错: %s", e)
		return e
	}

	log.Println(id)
	return e
}

func GetUserCredential(loginName string) (string, error) {
	stmt, err := dbConn.Prepare(`SELECT "password" from "user" where login_name = $1`)

	if err != nil {
		log.Printf("查询用户报错:%s", err)
		return "", err
	}
	defer stmt.Close()
	var pwd string
	stmt.QueryRow(loginName).Scan(&pwd)

	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmt, err := dbConn.Prepare(`DELETE FROM "user" where login_name = $1 and "password" = $2 `)
	if err != nil {
		log.Printf("删除用户报错:%s", err)
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(loginName, pwd)

	log.Println(result)
	return err
}

func AddVideo(authorId int, name string) (*defs.VideoInfo, error) {
	id := util.NewUUID()
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")
	sql := `insert into video_info (id,author_id,"name",display_ctime) values($1,$2,$3,$4)`

	stmt, e := dbConn.Prepare(sql)
	if e != nil {
		return nil, e
	}

	defer stmt.Close()
	_, e = stmt.Exec(id, authorId, name, ctime)
	if e != nil {
		return nil, e
	}

	return &defs.VideoInfo{Id: id, AuthorId: authorId, Name: name, DisplayCtime: ctime}, nil

}

func AddComment(authorId int,videoId string,content string) error {
	sql :=  `insert into comment (id,author_id,video_id,content) VALUES ($1,$2,$3,$4)`
	stmt, e := dbConn.Prepare(sql)
	if e != nil{
		log.Printf("添加评论出错:%s",e)
		return e
	}

	defer stmt.Close()

	_, e = stmt.Exec(util.NewUUID(), authorId, videoId, content)

	return e
}
