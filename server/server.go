package server

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gorpc/message"
	"gorpc/server/node"
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/serviceProto/protocol/protocol"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type ClientConnInfo struct {
	timeout  int
	retry    int
	maxRetry int
	conn     *io.ReadWriteCloser
}

type Server struct {
	handler   *serviceHandler.ServiceHandler
	connInfos sync.Map
}

func (s Server) Accept(lis *net.Listener) {
	for {
		conn, err := (*lis).Accept()
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		go s.handleConn(&conn)
	}
}

func (s Server) GetHandle() *serviceHandler.ServiceHandler {
	return s.handler
}

func (s Server) handleConn(_conn *net.Conn) {
	var conn io.ReadWriteCloser = *_conn

	defer conn.Close()

	s.connInfos.LoadOrStore(_conn, ClientConnInfo{
		500,
		0,
		5,
		&conn,
	})

	// 0 - 1: encode_way
	// 2 - 3: compress_way
	// 4 - 7: header length
	// ...header
	// n - n+3: body length
	// ...body
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
	var codec message.ServerCodec

	if meta[0] == 'B' {
		//codec = message.NewGobCodec(conn)
	} else if meta[0] == 'A' {
		codec = message.NewTlvCodec(&conn)
	}

	rh := protocol.RPCResponseHeader{
		Type: "acceptHandShake",
	}
	data, err := proto.Marshal(&rh)
	if err != nil {
		fmt.Println("send acceptHandShake error:", err)
	}

	codec.WriteWithLength(data)

	s.handleMessage(&codec, _conn)
}

func (s Server) parseTLVProto(conn *io.ReadWriteCloser) (*service.ServiceCall, error) {
	tlver := message.NewTlvCodec(conn)
	return tlver.ParseRequest(s.handler)
}

func (s Server) handleMessage(codec *message.ServerCodec, _conn *net.Conn) {
	var conn io.ReadWriteCloser = *_conn
	sending := sync.Mutex{}

	for {
		call, err := s.parseTLVProto(&conn)
		if err != nil {
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
				infos, ok := s.connInfos.Load(_conn)
				if !ok {
					panic("load connInfos err:" + err.Error())
				}
				if info, ok := infos.(ClientConnInfo); !ok {
					panic("unexpect error")
				} else {
					if info.retry >= info.maxRetry {
						fmt.Println("disconnect addr:", (*_conn).RemoteAddr().String())
						return
					}
					info.retry++
					s.connInfos.Store(_conn, info)
					time.Sleep(time.Duration(info.timeout) * time.Millisecond)
				}
			}
			continue
		}

		go s.handle(
			call,
			codec,
			&sending,
		)
	}
}

func (s Server) SendError(call *service.ServiceCall, codec *message.ServerCodec, mu *sync.Mutex, err string) {
	mu.Lock()
	defer mu.Unlock()

	(*codec).Write(&message.ClientHeader{
		"serviceResponse",
		call.Seq,
		err,
		nil,
	})
}

func (s Server) handle(call *service.ServiceCall, codec *message.ServerCodec, mu *sync.Mutex) {
	s.handler.Call(call.Method, call.Args, call.Reply)

	mu.Lock()
	defer mu.Unlock()

	responseHeader := protocol.RPCResponseHeader{
		Type: "serviceResponse",
		Seq:  call.Seq,
	}

	if r, ok := call.Reply.(proto.Message); !ok {
		responseHeader.Error = "service error: reply not message"
	} else {
		if rd, err := proto.Marshal(r); err != nil {
			responseHeader.Error = err.Error()
		} else {
			responseHeader.Reply = rd
		}
	}
	(*codec).WriteHeader(&responseHeader)
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
		sync.Map{},
	}
	s.Accept(&lis)
}
