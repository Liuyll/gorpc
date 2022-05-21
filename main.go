package main

import (
	"encoding/gob"
	"fmt"
	"github.com/golang/protobuf/proto"
	client2 "gorpc/client"
	"gorpc/serviceProto/service/test/Add"
	"gorpc/test"
	"time"
)

func main() {
	gob.Register(test.Args{})

	client := client2.NewClient()

	args, err := proto.Marshal(&Add.Payload{First: 1, Second: 2})
	if err != nil {
		panic(err)
	}

	reply := new(Add.Result)
	if err := client.CallWithTlv("test.AddProto", args, reply); err != nil {
		panic(err)
	} else {
		fmt.Println("RET:", reply.Result)
	}

	//for i := 0; i < 5; i++ {
	//	go func(i int) {
	//		ret := new(int)
	//		if err := client.Call("test.Add", test.Args{
	//			First:  i,
	//			Second: i,
	//		}, ret); err != nil {
	//			fmt.Println("err:", err)
	//		} else {
	//			fmt.Println("call end:", *ret, " i:", i)
	//		}
	//	}(i)
	//}

	time.Sleep(time.Duration(50) * time.Second)
}
