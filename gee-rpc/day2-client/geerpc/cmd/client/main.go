package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"geerpc"
	"geerpc/client"
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

			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}

			log.Println("reply:", reply)
		}(i)
	}

	wg.Wait()
}
