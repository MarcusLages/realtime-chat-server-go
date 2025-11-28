package main

import (
	"bufio"
	"net"
)

type ProxyWorker struct {
	conn        net.Conn
	chat_server *ChatServer
	user        User
	reader_buf  *bufio.Reader
	writer_buf  *bufio.Writer
}

const UserBufSize int = 15

func New_proxy_worker(conn net.Conn, chat_server *ChatServer) ProxyWorker {
	return ProxyWorker{
		conn:        conn,
		reader_buf:  bufio.NewReader(conn),
		writer_buf:  bufio.NewWriter(conn),
		chat_server: chat_server,
		user:        New_user(UserBufSize),
	}
}

func (p *ProxyWorker) Start() {

}
