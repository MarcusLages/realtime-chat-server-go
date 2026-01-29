# Realtime Chat Server (Go)
A real-time chat application built with Go, supporting one chat rooms and multiple concurrent users.
## Overview
This is a raw TCP socket-based chat server implementation in Go that enables real-time communication between multiple clients in the same chat room.
## Project Structure
```bash
realtime-chat-server-go/
├── client/            # Client-side shell code
├── echo/              # Simple example of an echo server (for initial shell testing)
├── server/            # Main server implementation
└── go.mod             # Go module dependencies
```
## Running the Project
1. Download Go dependencies.
```bash
go mod download
```
2. Start the server. The server will start listening on the default port (typically `localhost:6666`, so no dynamic change of the server port).
```bash
go run server/server.go
```
3. Start the client shell on a separate terminal or use any available pure TCP shell that uses line based messages by connecting to port `6666`
```bash
cd client
go run client.go
```
## Available Commands
- All of the commands are case insensitive.
- Users are automatically logged out once their shell/connection is closed.
##### `/NCK <nickname>` 
Login with a username so you can send/receive messages.
##### `/LST` 
Lists logged users.
##### `/MSG <recipients> <message>` 
Sends the same message to all recipients. You can have multiple target recipients by separating them by comma (no spaces, just comma)
- Ex: `/MSG user1,user2 Hello everyone.`
##### `/GRP <groupname> <users>` 
Creates a group with all the users. Whenever a message is sent to a group, the message will be broadcasted to all the users in that group. The group name must start with hash (`#`). To add multiple users, separate them by comma similarly to the `/MSG` command.
## Extra Notes
Access the Elixir version: [MarcusLages/realtime-chat-server](https://github.com/MarcusLages/realtime-chat-server)

