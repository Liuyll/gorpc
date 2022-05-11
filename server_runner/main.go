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

	server.StartServer(":8765", &testHandler)
}
