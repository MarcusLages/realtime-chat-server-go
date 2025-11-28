package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type ProxyWorker struct {
	conn        net.Conn
	chat_server *ChatServer
	user        User
	reader_buf  *bufio.Reader       // IO Buffer to read from socket
	writer_buf  *bufio.Writer       // IO Buffer to write to socket
	group_name  map[string][]string // Group name -> list of users
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

// From ChatServer to ProxyWorker
func (p *ProxyWorker) process_server_res() {
	for res := range p.user.res_chn {
		var output string
		if res.From != ServerName {
			output += res.From
		}
		// New line is very important since most sockets read input
		// from line to line
		output += res.Data + "\n"

		if _, err := p.writer_buf.WriteString(output); err != nil {
			log.Printf("Writing error to socket: %v", err)
		}

		p.writer_buf.Flush()
	}
}

// From ProxyWorker to ChatServer
func (p *ProxyWorker) process_cmd(line string) {
	tok := strings.Fields(strings.TrimSpace(line))
	if len(line) == 0 || len(tok) == 0 {
		return
	}

	cmd := strings.ToUpper(tok[0])
	switch Cmd(cmd) {
	case NCK:
		nick := tok[1]
		req := Request{p.user, NCK, []string{nick}} // Data is of slice type
		p.chat_server.Send_request(req)
	case LST:
		req := Request{p.user, LST, []string{}}
		p.chat_server.Send_request(req)
	case MSG:
		// p.handle_msg()
	case GRP:
		// p.handle_grp()
	default:
		err_msg := fmt.Sprintf("Unknown command: %s\n", cmd)
		p.writer_buf.WriteString(err_msg)
	}
}

func (p *ProxyWorker) handle_nck() {

}
