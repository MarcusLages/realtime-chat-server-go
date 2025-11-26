// echo client
package main

import (
  "bufio"
  "fmt"
  "log"
  "net"
  "os"
)

func main() {
  conn, err := net.Dial("tcp", "localhost:6666")
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  stdinScanner := bufio.NewScanner(os.Stdin)
  connScanner := bufio.NewScanner(conn)

  for {
    fmt.Printf("> ")
    if !stdinScanner.Scan() {
      break
    }
    fmt.Fprintf(conn, "%s\n", stdinScanner.Text())
    if !connScanner.Scan() {
      break
    }
    fmt.Println(connScanner.Text())
  }
}
