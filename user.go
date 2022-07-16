package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userIp := conn.RemoteAddr().String()
	user := &User{
		C:    make(chan string),
		Name: userIp,
		Addr: userIp,
		Conn: conn,
	}
	go user.ListenUserMessage()
	return user
}

func (t *User) ListenUserMessage() {
	for {
		msg := <-t.C
		t.Conn.Write([]byte(msg + "\n"))
	}
}
