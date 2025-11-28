package chat_server

import "regexp"

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

// TODO: Do not forget to check for renaming
func (s *ChatServer) process_nck(req Request) {
	nick := req.Data
	if !is_nick_valid(nick) {

	}

	res := Response{}
	req.From.Send_res(res)
}

func (s *ChatServer) process_lst(req Request) {

}

func (s *ChatServer) process_msg(req Request) {

}

func (s *ChatServer) process_logout(req Request) {

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
