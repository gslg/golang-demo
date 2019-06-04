package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main(){
	//1.连接远程tcp客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8888")

	if err != nil{
		log.Fatalf("connect error : %v",err)
		return
	}

	//2.写数据

	for {
		_, err := conn.Write([]byte("哈喽,world"))
		if err != nil{
			log.Fatalf("write conn error :%v",err)
			return
		}

		buf := make([]byte,512)
		cnt,err := conn.Read(buf)
		if err != nil{
			log.Fatalf("write conn error :%v",err)
			return
		}

		fmt.Printf("Received from server: %s,cnt=%d.\n",buf,cnt)

		time.Sleep(time.Second)
	}
}
