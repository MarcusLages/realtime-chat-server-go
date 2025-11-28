package main

import (
	"log"
	"net"
)

func main() {
	chat_server := ChatServer{}
	go chat_server.Start()

	addr := ":6666"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go start_new_worker(conn, &chat_server)
	}
}

func start_new_worker(conn net.Conn, chat_server *ChatServer) {
	log.Printf("New connection: %s", conn.RemoteAddr())
	worker := New_proxy_worker(conn, chat_server)
	go worker.Start()
}
