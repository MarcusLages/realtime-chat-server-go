package chat_server

type User struct {
	Nick    string
	res_chn chan Response // Channel to send responses back to user
}

func (u *User) Send_res(res Response) {
	u.res_chn <- res
}
