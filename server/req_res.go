package main

import "fmt"

// Name used when sending server responses (commonly error responses)
const ServerName string = "CHAT_SERVER"

type Request struct {
	From User
	Cmd  Cmd      // Sum type representing possible commands
	Data []string // Extra req data
}

// From/To are user nicks
type Response struct {
	From string // User nick, not UUID
	To   string // User nick, not UUID
	Data string
}

// Helper error response generators
// Only used in ChatServer
func Succ_server_res(dest, msg string) Response {
	return Response{ServerName, dest, msg}
}

func Err_res(dest, err_msg string) Response {
	return Response{ServerName, dest, err_msg}
}

func Err_nck_doesnt_exist(dest, non_exist_nick string) Response {
	err_msg := fmt.Sprintf("%s doesn't exist.", non_exist_nick)
	return Err_res(dest, err_msg)
}

func Err_invalid_cmd(dest string, inval_cmd Cmd) Response {
	err_msg := fmt.Sprintf("`%s` is an invalid command.", string(inval_cmd))
	return Err_res(dest, err_msg)
}

func Err_invalid_nick(dest string) Response {
	err_msg := "Invalid nickname. Must start with a letter, include only alphanumeric chars and underscores, and have a max length of 10."
	return Err_res(dest, err_msg)

}

func Err_nick_already_exists(dest string) Response {
	err_msg := "Invalid nickname. Nickname already exists."
	return Err_res(dest, err_msg)

}

func Err_unauthorized(dest string) Response {
	err_msg := "Unauthorized. You must have a NCK first using the `/NCK` cmd."
	return Err_res(dest, err_msg)
}
