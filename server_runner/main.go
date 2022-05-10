package main

import (
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/test"
)

func main() {
	testService := service.NewService(test.TestService{})
	testHandler := serviceHandler.NewServiceHandler()
	testHandler.Register("test", testService)
}
