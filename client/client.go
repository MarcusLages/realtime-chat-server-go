package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	port := 6666
	conn_url := fmt.Sprintf("localhost:%d", port)

	conn, err := net.Dial("tcp", conn_url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	done := make(chan bool)
	key_scan := bufio.NewScanner(os.Stdin)
	sock_scan := bufio.NewScanner(conn)

	go read_from_keyboard(key_scan, conn, done)
	go read_from_socket(sock_scan, done)

	// Just need to wait for either routine to be finished to finish the program
	<-done

	fmt.Println("Connection closed")
}

// Coroutine to print/read from keyboard
func read_from_keyboard(key_scan *bufio.Scanner, conn net.Conn, done chan bool) {
	defer func() { done <- true }()
	for {
		fmt.Print("> ")
		if !key_scan.Scan() {
			break
		}
		fmt.Fprintf(conn, "%s\n", key_scan.Text())
	}
}

// Coroutine to print/read from server socket
func read_from_socket(sock_scan *bufio.Scanner, done chan bool) {
	defer func() { done <- true }()
	for {
		if !sock_scan.Scan() {
			break
		}
		// This unicode string clears the line ("\r\033[K") to clear up user input
		// and the arrow prefix, writes down the socket response and then,
		// prints "> " on a new line
		fmt.Printf("\r\033[K%s\n> ", sock_scan.Text())
	}
}
