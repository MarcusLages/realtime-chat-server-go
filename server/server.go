package chat_server

// Cmd sum type with the possible commands
type Cmd string

const (
	NCK    Cmd = "/NCK"
	MSG    Cmd = "/MSG"
	LST    Cmd = "/LST"
	LOGOUT Cmd = "/LOGOUT"
)

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
		res := s.process_nck(req)
		s.send_to(res.To, res)
	case LST:
		res := s.process_lst(req)
		s.send_to(res.To, res)
	case MSG:
		res := s.process_msg(req)
		s.send_to(res.To, res)
	case LOGOUT:
		res := s.process_logout(req)
		s.send_to(res.To, res)
	default:
		err_res := Err_invalid_cmd(req.From, req.Cmd)
		s.send_to(err_res.To, err_res)
	}
}

// TODO: Do not forget to check for renaming
func (s *ChatServer) process_nck(req Request) Response {
	return Response{}
}

func (s *ChatServer) process_lst(req Request) Response {
	return Response{}
}

func (s *ChatServer) process_msg(req Request) Response {
	return Response{}
}

func (s *ChatServer) process_logout(req Request) Response {
	return Response{}
}

func (s *ChatServer) send_to(dest_nick string, res Response) {
	dest, exists := s.users[dest_nick]
	src := s.users[res.From]

	if !exists {
		err_res := Err_nck_doesnt_exist(res.From, dest_nick)
		src.Send_res(err_res)
	} else {
		dest.Send_res(res)
	}
}
