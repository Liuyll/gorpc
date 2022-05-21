package message

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"gorpc/service"
	"gorpc/serviceHandler"
	"io"
	"runtime/debug"
	"sync"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	dec  *gob.Decoder
	enc  *gob.Encoder
	buf  *bufio.Writer
	mu   *sync.Mutex
}

func NewGobCodec(conn io.ReadWriteCloser) GobCodec {
	buf := bufio.NewWriter(conn)
	return GobCodec{
		conn,
		gob.NewDecoder(conn),
		gob.NewEncoder(buf),
		buf,
		new(sync.Mutex),
	}
}

func (c GobCodec) Write(h interface{}, body interface{}) {
	defer func() {
		c.mu.Unlock()
		if err := c.buf.Flush(); err != nil {
			fmt.Println("flush error:", err)
		}
	}()

	c.mu.Lock()
	if err := c.enc.Encode(h); err != nil {
		fmt.Println("encode header err: \n", err, fmt.Sprintf("debug stack: %s", debug.Stack()))
	}
	if body != nil {
		if err := c.enc.Encode(body); err != nil {
			fmt.Println("encode body err: \n", err, fmt.Sprintf("debug stack: %s", debug.Stack()))
		}
	}
}

func (c GobCodec) ReadHeader(h interface{}) error {
	if err := c.dec.Decode(h); err != nil {
		return err
	}

	return nil
}

func (c GobCodec) ReadBody(body interface{}) error {
	return c.ReadHeader(body)
}

func (c GobCodec) Close() {
	c.conn.Close()
}

func (c GobCodec) ParseRequest(handler *serviceHandler.ServiceHandler) (*service.ServiceCall, error) {
	var h = RPCHeader{}

	if err := c.ReadHeader(h); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}

	var call service.ServiceCall
	serviceMethod := h.ServiceMethod
	if err, method := handler.ResolveServiceMethod(serviceMethod); err != nil {
		return nil, errors.New("ResolveServiceMethod err: " + err.Error())
	} else {
		fmt.Println("start handle service:", h.ServiceMethod)

		var args = method.NewArgs()
		body := RPCBody{
			args,
		}

		if err := c.ReadBody(&body); err != nil {
			if err != io.EOF {
				fmt.Println("read body err:", err)

				call = service.NewServiceCall(nil, nil, nil, h.Seq)
				return &call, err
			}
		}

		call = service.NewServiceCall(method, body.Args, method.NewReply(), h.Seq)
		return &call, nil
	}


	return &call, nil
}
