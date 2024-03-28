package main

import (
	"fmt"
	"log"
	"syscall"
)

func client() int {
	//creating a tcp socket
	sockfd, err := syscall.Socket(syscall.AF_INET,
		syscall.SOCK_STREAM,
		0)
	if err != nil {
		log.Fatalf("Socket creation failed with error: %v", err)
		return 1
	}
	defer syscall.Close(sockfd)

	//setting socket options
	val := 1
	if err := syscall.SetsockoptInt(
		sockfd,
		syscall.SOL_SOCKET,
		syscall.SO_REUSEADDR,
		val); err != nil {
		log.Fatalf("Failed to set socket options: %v", err)
		return 1
	}

	var addr syscall.SockaddrInet4
	addr.Port = 1234
	addr.Addr = [4]byte{127, 0, 0, 1} //loopbackaddress

	fmt.Println("Socket bound successfully to address")

	//try to connect to the server
	if err := syscall.Connect(sockfd, &addr); err != nil {
		log.Fatalf("Error when trying to connect with server: %v", &err)
	}

	//writing to the buffer of the server
	message := []byte("Hello")
	if _, err = syscall.Write(sockfd, message); err != nil {
		log.Fatalf("Write failed with error: %v", err)
	}
	//read from server
	buffer := make([]byte, 64)
	if _, err := syscall.Read(sockfd, buffer); err != nil {
		log.Fatalf("Read failed with error: %v", err)
		return 1
	}
	fmt.Printf("server says: %s\n", string(buffer))
	defer syscall.Close(sockfd)
	return 0
}
