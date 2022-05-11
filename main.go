package main

import (
	"encoding/gob"
	"fmt"
	client2 "gorpc/client"
	"gorpc/test"
	"time"
)

func main() {
	gob.Register(test.Args{})

	client := client2.NewClient()

	for i := 0; i < 5; i++ {
		go func(i int) {
			ret := new(int)
			if err := client.Call("test.Add", test.Args{
				First:  i,
				Second: i,
			}, ret); err != nil {
				fmt.Println("err:", err)
			} else {
				fmt.Println("call end:", *ret, " i:", i)
			}
		}(i)
	}

	time.Sleep(time.Duration(5) * time.Second)
}
