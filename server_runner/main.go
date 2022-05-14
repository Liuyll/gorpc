package main

import (
	"gorpc/server"
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/test"
)

func main() {
	//gob.Register(test.Args{})

	testService := service.NewService(test.TestService{})
	testHandler := serviceHandler.NewServiceHandler()
	testHandler.Register("test", testService)

	go server.StartServer("testservice", ":8765", &testHandler)
	go server.StartServer("testservice", ":8766", &testHandler)

	select {

	}
}
