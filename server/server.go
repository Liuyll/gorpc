package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
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
	buf := make([]byte, 12)

	reader := bufio.NewReader(conn)
	for {
		var n int
		if _n, err := reader.Read(buf); err != nil {
			n = _n
			if err != io.EOF {
				fmt.Println("read tcp body err:", err)
				return
			}
		}
		fmt.Println(string(buf[:n]))
	}
}

func Accept(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	Server{}.Accept(&lis)
}
