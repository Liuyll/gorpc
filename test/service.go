package test

import "fmt"

type TestService struct {

}

type Args struct {
	First int
	Second int
}

func (t TestService) Add(args Args) int {
	ret := args.First + args.Second
	fmt.Println("return ret:", ret)
	return ret
}