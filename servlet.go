package main

/*
梳理一下下面的代码：
首先新建一个server
然后启动server
	监听端口，接受连接，获取conn
	一个新的goroutine，负责监听BroadcastMessage
	然后启动接收用户消息后的处理函数
ListenMessage：监听BroadcastMessage，一但收到，给所有用户的C发送消息，user的C是一个chan string，收到消息后，发送给客户端
handleRequest：添加用户到在线用户列表，并给BroadcastMessage发送消息
*/

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip                string
	Port              int
	BroadcastMessage  chan string
	OlineUsersMap     map[string]*User
	OlineUsersMapLock sync.Mutex
}

func NewServer(Ip string, Port int) *Server {
	server := &Server{
		Ip:               Ip,
		Port:             Port,
		BroadcastMessage: make(chan string),
		OlineUsersMap:    make(map[string]*User),
	}
	return server
}

func (t *Server) StartServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.Ip, t.Port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	go t.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		go t.handleRequest(conn)
	}

}

func (t *Server) handleRequest(conn net.Conn) {
	user := NewUser(conn)

	t.OlineUsersMapLock.Lock()
	t.OlineUsersMap[user.Name] = user
	t.OlineUsersMapLock.Unlock()

	t.BroadcastMessage <- user.Name + " login"

}

func (t *Server) ListenMessage() {
	for {
		msg := <-t.BroadcastMessage
		t.OlineUsersMapLock.Lock()
		for _, user := range t.OlineUsersMap {
			user.C <- msg
		}
		t.OlineUsersMapLock.Unlock()
	}
}
