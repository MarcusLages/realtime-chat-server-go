package main

import (
	"crypto/rand"
	"fmt"
)

type User struct {
	// Unique random UUID used to identify ID
	ID      string        // I chose to do it randomly to avoid global vars and race conditions
	res_chn chan Response // Buffered channel to send responses back to user
}

// User uses a buffered channed to avoid blocking the client
func New_user(buff_size int) User {
	return User{
		ID:      gen_id(),
		res_chn: make(chan Response, buff_size),
	}
}

// User is responsible for parsing and formatting the Response
func (u *User) Send_res(res Response) {
	u.res_chn <- res
}

// Generate UUID
func gen_id() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
