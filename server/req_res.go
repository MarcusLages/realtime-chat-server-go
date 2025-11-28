package chat_server

import "fmt"

// Name used when sending server responses (commonly error responses)
const ServerName string = "CHAT_SERVER"

// From is a user nick
type Request struct {
	From string
	Cmd  Cmd    // Sum type representing possible commands
	Data string // Extra req data
}

// From/To are user nicks
type Response struct {
	From string
	To   string
	Data string
}

// Helper error response generators
func Err_res(dest, err_msg string) Response {
	return Response{ServerName, dest, err_msg}
}

func Err_nck_doesnt_exist(dest, non_exist_nick string) Response {
	err_msg := fmt.Sprintf("%s doesn't exist.", non_exist_nick)
	return Err_res(dest, err_msg)
}

func Err_invalid_cmd(dest, inval_cmd string) Response {
	err_msg := fmt.Sprintf("`%s` is an invalid command.", inval_cmd)
	return Err_res(dest, err_msg)
}

func Err_unauthorized(dest string) Response {
	err_msg := "Unauthorized. You must have a NCK first using the `/NCK` cmd."
	return Err_res(dest, err_msg)
}
