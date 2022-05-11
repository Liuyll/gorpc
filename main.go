package main

import (
	"encoding/gob"
	"fmt"
	client2 "gorpc/client"
	"gorpc/server"
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/test"
	"time"
)

func createServer(done2 chan byte) {
	testService := service.NewService(test.TestService{})
	testHandler := serviceHandler.NewServiceHandler()
	testHandler.Register("test", testService)

	done := server.Accept(":8765", &testHandler)
	_ = <-done
	done2 <- 'A'
}

func main() {
	gob.Register(test.Args{})

	createDone := make(chan byte)
	fmt.Println("createServer")
	go createServer(createDone)

	_ = <-createDone

	client := client2.NewClient()

	for i := 0; i < 5; i++ {
		ret := new(int)
		go client.Call("test.Add", test.Args{
			First:  i,
			Second: i,
		}, ret)
	}

	time.Sleep(time.Duration(5) * time.Second)
}
