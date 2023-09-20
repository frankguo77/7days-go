package main

import (
	"net"
	"log"
	"geerpc"
	"geerpc/server"
)

// func startServer(addr chan string) {
// 	l, err := net.Listen("tcp", ":0")

// 	if err != nil {
// 		log.Fatal("network error:", err)
// 	}

// 	log.Println("start rpc server on ", l.Addr())
// 	addr <- l.Addr().String()
// 	server.Accept(l)
// }

func startServer(addr string) {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on ", l.Addr())
	// addr <- l.Addr().String()
	server.Accept(l)
}

func main() {
	log.SetFlags(0)

	// addr := make(chan string)
	startServer(geerpc.ADDR)
}