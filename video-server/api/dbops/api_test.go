package dbops

import (
	"log"
	"testing"
)

func clearTables() {
	dbConn.Exec("truncate user")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUserCredential)
	t.Run("Get", testGetUserCredential)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUserCredential(t *testing.T) {
	err := AddUserCredential("liuguo", "123456")
	if err != nil {
		t.Errorf("Error Of AddUser:%v", err)
	}
}

func testGetUserCredential(t *testing.T) {
	pwd, err := GetUserCredential("liuguo")
	if err != nil {
		t.Errorf("Error Of GetUser:%v", err)
	}

	if pwd != "123456" {
		t.Error("User pwd Error")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("liuguo", "123456")
	if err != nil {
		t.Errorf("Error Of DeleteUser:%v", err)
	}
}

func testRegetUser(t *testing.T) {
	pwd, err := GetUserCredential("liuguo")
	if err != nil {
		t.Errorf("Error Of RegetUser:%v", err)
	}
	if pwd != "" {
		t.Error("Delete User Failed")
	}
}

func TestAddVideo(t *testing.T) {
	dbConn.Exec("truncate video_info")
	v, err := AddVideo(123, "我是一个视频")

	if err != nil {
		t.Errorf("添加视频出错:%s", err)
	}

	if v == nil {
		t.Error("添加视频出错")
	}

	log.Println(*v)

	dbConn.Exec("truncate video_info")

}

func TestAddComment(t *testing.T) {
	dbConn.Exec("truncate comment")
	e := AddComment(123, "223", "这是一个评论")
	if e != nil {
		t.Errorf("添加评论出错:%s",e)
	}
	t.Log("添加评论成功")

	dbConn.Exec("truncate comment")
}
