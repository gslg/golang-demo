package server

import (
	"fmt"
	"log"
	"net"
)

type IServer interface {
	//启动
	Start()
	//停止方法
	Stop()
	//服务方法
	Serve()
}

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	//1.创建一个套接字
	log.Printf("Try Start Server [%s] at IP:%s,Port:%d.\n", s.Name, s.IP, s.Port)
	go func() {
		//2.监听该socket
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			log.Fatalf("resolve tcp address error:%v", err)
			return
		}

		//3.阻塞的等待连接，处理业务

		listener, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			log.Fatalf("Listen tcp at address error:%v", err)
			return
		}

		log.Printf("Start Server[%s] success,listen at [%d].....",s.Name,s.Port)

		for {
			//阻塞直到客户端连接
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Printf("Accept error:%v", err)
				continue
			}

			log.Printf("Accept from client:%v",conn.RemoteAddr())

			//处理业务逻辑
			go func() {
				for {
					buf := make([]byte,512)
					cnt,err := conn.Read(buf)
					if err != nil{
						log.Printf("read from client error:%v", err)
						break
					}
					log.Printf("Received from client: %s,cnt=%d",buf,cnt)
					if _,err = conn.Write(buf[:cnt]);err != nil {
						log.Printf("write to client error:%v", err)
						break
					}
				}

			}()

		}
	}()
}

func (s *Server) Stop() {
  //停止服务器，回收资源等
}

func (s *Server) Serve() {
	s.Start()

	//这里可以做一些其他的事
	//阻塞
	select {

	}
}

//实例化一个server服务器
func New(name string) IServer {
	s := Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
	}

	return &s
}
