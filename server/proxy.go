package main

import (
	"bufio"
	"log"
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
	defer p.cleanup()

	go p.process_server_res()

	// Client socket reading loop
	for {
		line, err := p.reader_buf.ReadString('\n')
		if err != nil { // nil represents EOF
			// Connection closed or error
			return
		}
		p.process_cmd(line)
	}
}

// Should be called after the user logs out, for clean up
// (done after socket is closed in this case)
// Logs out user from ChatServer and closes the ProxyWorker connection
func (p *ProxyWorker) cleanup() {
	logout_req := Request{
		From: p.user,
		Cmd:  LOGOUT,
	}
	p.chat_server.Send_request(logout_req)
	p.conn.Close()
	log.Printf("Connection closed: %s", p.conn.RemoteAddr())
}

// From GameServer to ProxyWorker
func (p *ProxyWorker) process_server_res() {

}

// From ProxyWorker to GameServer
func (p *ProxyWorker) process_cmd(line string) {

}
