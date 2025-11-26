// echo server
package main

import (
  "bufio"
  "fmt"
  "log"
  "net"
)

func main() {
  ln, err := net.Listen("tcp", ":6666")
  if err != nil {
    log.Fatal(err)
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Fatal(err)
    }
    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  defer conn.Close()

  scanner := bufio.NewScanner(conn)

  for scanner.Scan() {
    fmt.Fprintf(conn, "%s\n", scanner.Text())
  }

  fmt.Println("connection closed:", conn)
}
