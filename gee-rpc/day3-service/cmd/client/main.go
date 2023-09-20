package main

import (
	"log"
	"sync"
	"time"
	"geerpc"
	"geerpc/client"
	"example/cmd"
)

func main() {
	client, _ := client.Dial("tcp", geerpc.ADDR)
	defer func() { client.Close() }()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			args := &cmd.Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}

			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}

	wg.Wait()
}
