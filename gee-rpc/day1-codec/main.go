package main

import (
	"geerpc/codec"
	"encoding/json"
	"fmt"
	"geerpc/server"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")

	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on ", l.Addr())
	addr <- l.Addr().String()
	server.Accept(l)
}

func main() {
	log.SetFlags(0)

	addr := make(chan string)
	go startServer(addr)


	conn, _ := net.Dial("tcp", <-addr)
	defer func(){conn.Close()}()

	time.Sleep(time.Second)
	json.NewEncoder(conn).Encode(server.DefaultOption)

	cc := codec.NewGobCodec(conn)

	for i:= 0; i < 5; i++ {
		h := &codec.Header {
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}

		cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		cc.ReadHeader(h)
		var reply string
		cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}