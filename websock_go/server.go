
package main

import (
	"fmt"
	"net"
)

func msg(msg string) {
	fmt.Println(msg)
}

func die(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
	panic(msg)
}

func doSomething(conn net.Conn) {
	rbuf := make([]byte, 64)
	n, err := conn.Read(rbuf)
	if err != nil {
		msg("read() error")
		return
	}
	fmt.Printf("client says: %s\n", string(rbuf[:n]))

	wbuf := []byte("world")
	_, err = conn.Write(wbuf)
	if err != nil {
		msg("write() error")
	}
}

func main() {
	// Create a TCP socket
	fd, err := net.Listen("tcp", ":1234")
	if err != nil {
		die("socket() failed")
	}
	defer fd.Close()

	// Accept connections in a loop
	for {
		conn, err := fd.Accept()
		if err != nil {
			msg("accept() error")
			continue
		}

		// Handle connection in a separate goroutine
		go func(conn net.Conn) {
			defer conn.Close()
			doSomething(conn)
		}(conn)
	}
}
