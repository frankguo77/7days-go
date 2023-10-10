package main

import (
	"net"
	"log"
	"geerpc"
	"geerpc/cmd"
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

type Foo int

type Args struct{Num1, Num2 int}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}


func startServer(addr string) {
	var foo cmd.Foo
	if err := geerpc.Server.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on ", l.Addr())
	// addr <- l.Addr().String()
	geerpc.server.Accept(l)
}

func main() {
	log.SetFlags(0)

	// addr := make(chan string)
	startServer(geerpc.ADDR)
}