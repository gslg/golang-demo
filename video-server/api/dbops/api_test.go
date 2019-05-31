package dbops

import "testing"

func TestMain(m *testing.M){

}

func TestUserWorkFlow(t *testing.T){

}

func testAddUserCridential(t *testing.T) {
	err := AddUserCridential("liuguo","123456")
	if err != nil{
		t.Errorf("Error:%v",err)
	}
}

func testGetUserCridential(t *testing.T) {

}

func testDeleteUser(t *testing.T) {

}
