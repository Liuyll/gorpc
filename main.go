package main

import (
	"bufio"
	"fmt"
	"gorpc/message"
	"gorpc/server"
	"net"
	"time"
)

func createServer(done2 chan byte) {
	done := server.Accept(":8765")
	_ = <-done
	done2 <- 'A'
}

func main() {
	createDone := make(chan byte)
	fmt.Println("createServer")
	go createServer(createDone)

	_ = <-createDone

	fmt.Println("connection")
	conn, err := net.Dial("tcp", "127.0.0.1:8765")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}

	writer := bufio.NewWriter(conn)
	if _, err := writer.Write([]byte{'A', 'B', 'C', 'D', '\n'}); err != nil {
		fmt.Println("write err:", err)
		return
	}

	writer.Flush()
	time.Sleep(time.Duration(1) * time.Second)

	codec := message.NewGobCodec(conn)
	codec.Write(&message.RPCHeader{
		ServiceMethod: "test",
	}, "client")

	h := &message.RPCHeader{}
	if err := codec.ReadHeader(h); err != nil {
		fmt.Println("read server header err:", err)
		return
	}
	fmt.Println("read service header:", h.ServiceMethod)

	time.Sleep(time.Duration(5) * time.Second)

}
