package server

import (
	"bufio"
	"fmt"
	"gorpc/message"
	"gorpc/serviceHandler"
	"gorpc/test"
	"io"
	"net"
)

type Server struct {
	handler *serviceHandler.ServiceHandler
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

	meta := buf
	var codec message.Codec

	if meta[0] == 'A' {
		codec = message.NewGobCodec(conn)
	}
	codec.Write(message.NewShakeClientHeader(), nil)

	s.handleMessage(&codec)
}

func (s Server) handleMessage(codec *message.Codec) {
	var h = &message.RPCHeader{}

	for {
		if err := (*codec).ReadHeader(h); err != nil {
			if err != io.EOF {
				fmt.Println("read header err:", err)
				(*codec).Close()
				return
			}
		}

		serviceMethod := h.ServiceMethod
		if err, method := s.handler.ResolveServiceMethod(serviceMethod); err != nil {
			fmt.Println("ResolveServiceMethod err:", err)
		} else {
			fmt.Println("start handler service:", h.ServiceMethod)

			var args = method.NewArgs()
			var reply = method.NewReply()
			body := message.RPCBody{
				Args: &args,
				Reply: reply,
			}
			if v, ok := args.(test.Args); ok {
				fmt.Println("rrrrrrrr:", v.First)
				v.First = 2
			}

			fmt.Println("ggggggggg")
			if err := (*codec).ReadBody(&body); err != nil {
				if err != io.EOF {
					fmt.Println("read body err:", err)
					(*codec).Close()
					return
				}
			}

			call := serviceCall{
				method,
				args,
				reply,
			}

			go s.handle(
				&call,
				codec,
			)
		}
	}
}

func (s Server) handle(call *serviceCall, codec *message.Codec) {
	s.handler.Call(call.method, call.args, call.reply)

	(*codec).Write(&message.ClientHeader{
		"serviceResponse",
		0,
		nil,
		call.reply,
	}, nil)

	fmt.Println("exec end")

}

func Accept(port string, handler *serviceHandler.ServiceHandler) <-chan byte {
	done := make(chan byte, 1)
	startServer(port, handler, done)
	return done
}

func startServer(port string, serviceHandler *serviceHandler.ServiceHandler, done chan<- byte) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}

	done <- 'A'
	s := Server{
		serviceHandler,
	}
	go s.Accept(&lis)
}
