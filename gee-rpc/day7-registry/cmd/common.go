package cmd

import (
	// "log"
	"time"
)

type Foo int

type Args struct{Num1, Num2 int}

func (f Foo) Sum(args Args, reply *int) error {
	// log.Printf("Foo: Sum: %+v\n", args)
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1))
	*reply = args.Num1 + args.Num2
	return nil
}