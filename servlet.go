package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(Ip string, Port int) *Server {
	server := &Server{Ip: Ip, Port: Port}
	return server
}

func (t *Server) StartServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.Ip, t.Port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
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
	defer conn.Close()
	var buf [512]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}

}
