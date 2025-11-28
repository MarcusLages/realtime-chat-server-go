package chat_server

// Cmd sum type with the possible commands
type Cmd string

const (
	NCK Cmd = "/NCK"
	MSG Cmd = "/MSG"
	LST Cmd = "/LST"
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

}
