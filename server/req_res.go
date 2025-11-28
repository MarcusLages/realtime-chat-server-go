package chat_server

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
