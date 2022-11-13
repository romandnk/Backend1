package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			conn.Close()
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		_, err := io.WriteString(conn, "Server: "+input.Text())
		if err != nil {
			return
		}
	}
}
