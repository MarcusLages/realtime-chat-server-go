package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

// Default channel buffer size among all the users using the proxy
const UserBufSize int = 15

// Regex to check for valid group names
const GroupRegex string = "^#[a-zA-Z][a-zA-Z0-9_]{0,9}$"

// New group command implemented by the proxy
const GRP Cmd = "/GRP"

type ProxyWorker struct {
	conn        net.Conn
	chat_server *ChatServer
	user        User
	reader_buf  *bufio.Reader       // IO Buffer to read from socket
	writer_buf  *bufio.Writer       // IO Buffer to write to socket
	groups      map[string][]string // Group name -> list of users
}

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
		p.write_to_socket(output)
	}
}

// Msg needs to include a '\n'
func (p *ProxyWorker) write_to_socket(msg string) {
	if _, err := p.writer_buf.WriteString(msg); err != nil {
		log.Printf("Error while writing to socket: %v", err)
	}

	p.writer_buf.Flush()
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
		p.handle_msg(tok[1:])
	case GRP:
		p.handle_grp(tok[:1])
	default:
		err_msg := fmt.Sprintf("Unknown command: %s\n", cmd)
		p.writer_buf.WriteString(err_msg)
	}
}

func (p *ProxyWorker) handle_msg(args []string) {
	if len(args) < 2 {
		p.write_to_socket("Wrong usage of /MSG. Usage: /MSG <dest> <msg>")
		return
	}

	dest_str := args[1]
	dest := remove_duplicates(p.expand_dest(dest_str))

	msg := strings.Join(args[2:], " ")
	for _, d := range dest {
		req := Request{
			From: p.user,
			Cmd:  MSG,
			Data: []string{d, msg},
		}
		p.chat_server.Send_request(req)
	}
}

func (p *ProxyWorker) handle_grp(args []string) {
	if len(args) < 2 {
		p.write_to_socket("Wrong usage of /GRP. Usage: /GRP #group user,user,...")
		return
	}

	group_name := args[1]
	users_str := args[2]

	if !is_valid_group_name(group_name) {
		p.write_to_socket("Invalid group name. Must start with hash (#), include only alphanumeric chars and underscores, and have a max length of 11 (counting the hash).\n")
		return
	}

	users := strings.Split(users_str, ",")
	p.groups[group_name] = users

	msg := fmt.Sprintf("Group %s added\n", group_name)
	p.write_to_socket(msg)
}

func (p *ProxyWorker) expand_dest(dest_str string) []string {
	input_dests := strings.Split(dest_str, ",")
	expanded_dest := []string{}

	for _, d := range input_dests {
		// Group parsing
		if strings.HasPrefix(d, "#") {
			users, ok := p.groups[d]
			if !ok {
				err_msg := fmt.Sprintf("Group %s doesn't exist\n", d)
				p.write_to_socket(err_msg)
			} else {
				expanded_dest = append(expanded_dest, users...)
			}
		} else {
			// User name (probably, the ChatServer will be the one to give that answer)
			expanded_dest = append(expanded_dest, d)
		}
	}

	return expanded_dest
}

func remove_duplicates(arr []string) []string {
	map_set := map[string]int{}
	for _, x := range arr {
		map_set[x] = 0
	}

	result := []string{}
	for x := range map_set {
		result = append(result, x)
	}

	return result
}

func is_valid_group_name(group_name string) bool {
	if len(group_name) < 2 ||
		len(group_name) > 11 ||
		strings.HasPrefix(group_name, "#") {

		return false
	}

	matched, _ := regexp.MatchString(GroupRegex, group_name)
	return matched
}
