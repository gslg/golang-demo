package main

import "tcp-server/server"

func main(){
	//创建servr
	s := server.New("test")
	s.Serve()
}
