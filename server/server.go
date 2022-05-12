package server

import (
	"bufio"
	"fmt"
	"gorpc/message"
	"gorpc/server/node"
	"gorpc/serviceHandler"
	"io"
	"net"
	"sync"
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
	sending := sync.Mutex{}

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
			fmt.Println("start handle service:", h.ServiceMethod)

			var args = method.NewArgs()
			var reply = method.NewReply()
			body := message.RPCBody{
				args,
			}

			if err := (*codec).ReadBody(&body); err != nil {
				if err != io.EOF {
					fmt.Println("read body err:", err)
					call := serviceCall{
						nil,
						nil,
						nil,
						h.Seq,
					}
					s.SendError(&call, codec, &sending, err.Error())
					return
				}
			}

			call := serviceCall{
				method,
				body.Args,
				reply,
				h.Seq,
			}

			go s.handle(
				&call,
				codec,
				&sending,
			)
		}
	}
}

func (s Server) SendError(call *serviceCall, codec *message.Codec, mu *sync.Mutex, err string) {
	mu.Lock()
	defer mu.Unlock()

	(*codec).Write(&message.ClientHeader{
		"serviceResponse",
		call.Seq,
		err,
		nil,
	}, nil)
}

func (s Server) handle(call *serviceCall, codec *message.Codec, mu *sync.Mutex) {
	s.handler.Call(call.method, call.args, call.reply)

	mu.Lock()
	defer mu.Unlock()

	(*codec).Write(&message.ClientHeader{
		"serviceResponse",
		call.Seq,
		"",
		call.reply,
	}, nil)
}

func StartServer(name string, port string, handler *serviceHandler.ServiceHandler) <-chan byte {
	done := make(chan byte, 1)
	startServer(name, port, handler, done)
	return done
}

func startServer(name string, port string, serviceHandler *serviceHandler.ServiceHandler, done chan<- byte) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}

	node.NewServiceNode(lis.Addr().String(), name, 10)

	done <- 'A'
	s := Server{
		serviceHandler,
	}
	s.Accept(&lis)
}
