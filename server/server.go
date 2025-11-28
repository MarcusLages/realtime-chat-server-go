package main

import (
	"log"
	"net"
)

const ChatServerBufLimit int = 100

func main() {
	// Start the centralized ChatServer as a go routine
	// Using an
	chat_server := ChatServer{
		req_chn:    make(chan Request, ChatServerBufLimit),
		users:      make(map[string]User),
		user_nicks: make(map[string]string),
	}
	go chat_server.Start()

	addr := ":6666"
	ln, err := net.Listen("tcp", addr) // Listener socket
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Listening on %s", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go start_new_worker(conn, &chat_server) // Start proxy worker attached to this socket
	}
}

func start_new_worker(conn net.Conn, chat_server *ChatServer) {
	log.Printf("New connection: %s", conn.RemoteAddr())
	worker := New_proxy_worker(conn, chat_server)
	worker.Start()
}
