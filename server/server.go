package server

import (
	"bufio"
	"fmt"
	"gorpc/message"
	"io"
	"net"
	"sync"
)

type Server struct {
}

func (s Server) Accept(lis *net.Listener) {
	for {
		conn, err := (*lis).Accept()
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s Server) handleConn(conn io.ReadWriteCloser) {
	defer conn.Close()

	// 0 - 1: encode_way
	// 2 - 3: compress_way
	// 4 - 7: header length
	// 8 - 11: body length
	var buf []byte

	reader := bufio.NewReader(conn)
	fmt.Println("start read")
	if _buf, err := reader.ReadBytes('\n'); err != nil {
		if err != io.EOF {
			fmt.Println("read tcp body err:", err)
			return
		}
	} else {
		buf = _buf
	}

	fmt.Println("break")
	meta := buf
	var codec message.Codec

	if meta[0] == 'A' {
		codec = message.NewGobCodec(conn)
	}

	s.handleMessage(&codec)
}

func (s Server) handleMessage(codec *message.Codec) {
	sending := &sync.Mutex{}
	var h = &message.RPCHeader{}

	for {

		if err := (*codec).ReadHeader(h); err != nil {
			if err != io.EOF {
				fmt.Println("read header err:", err)
				(*codec).Close()
				return
			}
		}

		fmt.Println("read client header")

		var body = new(string)
		if err := (*codec).ReadBody(body); err != nil {
			if err != io.EOF {
				fmt.Println("read body err:", err)
				(*codec).Close()
				return
			}
		}

		go s.handle(
			&message.RPCMessage{
				h,
				&message.RPCBody{},
			},
			codec,
			sending,
		)
	}
}

func (s Server) handle(msg *message.RPCMessage, codec *message.Codec, sending *sync.Mutex) {
	defer sending.Unlock()

	sending.Lock()
	fmt.Println("call method", msg.H.ServiceMethod)
	(*codec).Write(&message.RPCHeader{
		ServiceMethod: msg.H.ServiceMethod,
	}, "qwe")
}

func Accept(port string) <-chan byte {
	done := make(chan byte, 1)
	startServer(port, done)
	return done
}

func startServer(port string, done chan<- byte) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}

	fmt.Println("lis create", lis.Addr().String())
	done <- 'A'
	go Server{}.Accept(&lis)
}
