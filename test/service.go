package test

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gorpc/serviceProto/service/test/Add"
)

type TestService struct {
}

type Args struct {
	First  int
	Second int
}

type Args2 struct {
	First1  string
	Second1 int
}

func (t TestService) Add(args Args) int {
	fmt.Println("call add")
	ret := args.First + args.Second
	return ret
}

func (t TestService) UnmarshalAdd(args []byte) (*Add.Payload, error) {
	ret := Add.Payload{}
	err := proto.Unmarshal(args, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (t TestService) AddProto(payload Add.Payload) Add.Result {
	ret := int(payload.First + payload.Second)
	fmt.Println("ret:", ret)
	return Add.Result{
		Result: int32(ret),
	}
}

func (t TestService) UnmarshalAddProto(args []byte) (*Add.Payload, error) {
	ret := Add.Payload{}
	err := proto.Unmarshal(args, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
