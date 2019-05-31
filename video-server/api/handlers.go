package main

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func CreateUser(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Create User Handler")
}

func Login(w http.ResponseWriter,r *http.Request){
	params := mux.Vars(r)
	username := params["username"]
	io.WriteString(w,"登陆用户: "+username)
}