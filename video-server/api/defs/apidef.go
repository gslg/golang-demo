package defs

type UserCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VideoInfo struct {
	Id string
	AuthorId int
	Name string
	DisplayCtime string
}
