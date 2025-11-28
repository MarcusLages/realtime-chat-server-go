package chat_server

import (
	"log"
	"regexp"
)

// Cmd sum type with the possible commands
type Cmd string

const (
	NCK    Cmd = "/NCK"
	MSG    Cmd = "/MSG"
	LST    Cmd = "/LST"
	LOGOUT Cmd = "/LOGOUT"
)

const NickRegex string = "^[a-zA-Z][a-zA-Z0-9_]{0,9}$"

type ChatServer struct {
	req_chn chan Request    // Request channel (receives commands)
	users   map[string]User // nicks -> User
}

func (s *ChatServer) Start() {
	for {
		req := <-s.req_chn
		s.process_command(req)
	}
}

func (s *ChatServer) Send_request(req Request) {
	s.req_chn <- req
}

func (s *ChatServer) process_command(req Request) {
	switch req.Cmd {
	case NCK:
		s.process_nck(req)
	case LST:
		s.process_lst(req)
	case MSG:
		s.process_msg(req)
	case LOGOUT:
		s.process_logout(req)
	default:
		err_res := Err_invalid_cmd(req.From.Nick, req.Cmd)
		req.From.Send_res(err_res)
	}
}

// */NCK
func (s *ChatServer) process_nck(req Request) {
	nick := req.Data[0]
	if !is_nick_valid(nick) {
		err_res := Err_invalid_nick(req.From.Nick)
		req.From.Send_res(err_res)
		return
	}

	ex_user, exists := s.users[nick]
	if exists && ex_user.Nick != req.From.Nick {
		err_res := Err_nick_already_exists(req.From.Nick)
		req.From.Send_res(err_res)
		return
	}

	// Renaming user
	if old_nick, exists := s.users[req.From.Nick]; exists {
		delete(s.users, req.From.Nick)
		old_nick.Nick = nick
		req.From.Nick = nick // For clarity
		s.users[nick] = old_nick
	} else {
		// New user
		req.From.Nick = nick // For clarity
		s.users[nick] = req.From
	}

	log.Printf("User '%s' was added to the chat.\n", nick)

	res := Succ_server_res(nick, "OK")
	req.From.Send_res(res)
}

// */LST
func (s *ChatServer) process_lst(req Request) {
	var user_list []string
	for nick := range s.users {
		user_list = append(user_list, nick)
	}

	// List formatting
	list_str := "["
	for i, nick := range user_list {
		if i > 0 && i < len(list_str) {
			list_str += ", "
		}
		list_str += nick
	}
	list_str += "]"

	res := Succ_server_res(req.From.Nick, list_str)
	req.From.Send_res(res)
}

// */MSG
func (s *ChatServer) process_msg(req Request) {
	dest_nick := req.Data[0]
	msg := req.Data[1]
	if _, authorized := s.users[req.From.Nick]; !authorized {
		err_res := Err_unauthorized(req.From.Nick)
		req.From.Send_res(err_res)
	}

	log.Printf("Sending msg from %s to %s.\n", req.From.Nick, msg)
	res := Response{
		From: req.From.Nick,
		To:   dest_nick,
		Data: msg,
	}
	s.send_msg(req.From, dest_nick, res)
}

// */LOGOUT
func (s *ChatServer) process_logout(req Request) {
	nick := req.From.Nick

	if _, exists := s.users[nick]; exists {
		delete(s.users, nick)
		log.Printf("Logging out user: %s\n", nick)
	}
}

func (s *ChatServer) send_msg(src User, dest_nick string, res Response) {
	dest, exists := s.users[dest_nick]

	if !exists {
		err_res := Err_nck_doesnt_exist(res.From, dest_nick)
		src.Send_res(err_res)
	} else {
		dest.Send_res(res)
	}
}

func is_nick_valid(nick string) bool {
	matched, _ := regexp.MatchString(NickRegex, nick)
	return matched
}
