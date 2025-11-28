package chat_server

import (
	"crypto/rand"
	"fmt"
)

type User struct {
	ID      string
	res_chn chan Response // Channel to send responses back to user
}

func New_user(buff_size int) User {
	return User{
		ID:      gen_id(),
		res_chn: make(chan Response, buff_size),
	}
}

func (u *User) Send_res(res Response) {
	u.res_chn <- res
}

// Generate UUID
func gen_id() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
