package main

import (
	"fmt"
	"log"
	"syscall"
)

func handleConnection(connfd int) {
	// Read from connection
	buffer := make([]byte, 64)
	n, err := syscall.Read(connfd, buffer)
	if err != nil {
		log.Fatalf("Read failed with error: %v", err)
		return
	}
	fmt.Printf("Client says: %s\n", string(buffer[:n]))

	// Write to connection
	message := []byte("world")
	_, err = syscall.Write(connfd, message)
	if err != nil {
		log.Fatalf("Write failed with error: %v", err)
		return
	}

	// Close connection
	syscall.Close(connfd)
}

func socket() {
	//creating a tcp socket
	sockfd, err := syscall.Socket(syscall.AF_INET,
		syscall.SOCK_STREAM,
		0)
	if err != nil {
		log.Fatalf("Socket creation failed with error: %v", err)
		return
	}
	defer syscall.Close(sockfd)

	//setting socket options
	val := 1
	if err := syscall.SetsockoptInt(sockfd,
		syscall.SOL_SOCKET,
		syscall.SO_REUSEADDR,
		val); err != nil {
		log.Fatalf("Failed to set socket options: %v", err)
		return
	}

	// Bind the socket to the wildcard address
	var addr syscall.SockaddrInet4
	addr.Port = 1234
	addr.Addr = [4]byte{0, 0, 0, 0} // Wildcard address

	if err := syscall.Bind(sockfd, &addr); err != nil {
		log.Fatalf("Bind failed with error: %v", err)
		return
	}

	fmt.Println("Socket bound successfully to address")

	//start listening for incoming connections
	if err := syscall.Listen(sockfd, syscall.SOMAXCONN); err != nil {
		log.Fatalf("Listen failed with error: %v", &err)
		return
	}

	for {
		// Accept incoming connections
		connfd, _, err := syscall.Accept(sockfd)
		if err != nil {
			log.Fatalf("Accept failed with error: %v", err)
			continue
		}

		// Handle each connection concurrently
		go handleConnection(connfd)
	}
}

//for how to setsockopt on golang I referred to
//https://iximiuz.com/en/posts/go-net-http-setsockopt-example
//and golang docs on syscall
