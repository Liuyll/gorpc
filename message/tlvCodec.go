package message

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/serviceProto/protocol/protocol"
	"gorpc/utils"
	"io"
	"sync"
	"time"
)

type TlvCodec struct {
	conn *io.ReadWriteCloser
	buf  *bufio.Writer
	mu   *sync.Mutex
}

type meta struct {
	encoding int
	compress int
}

func NewTlvCodec(conn *io.ReadWriteCloser) *TlvCodec {
	buf := bufio.NewWriter(*conn)
	mu := new(sync.Mutex)
	return &TlvCodec{
		conn,
		buf,
		mu,
	}
}

func (c TlvCodec) WriteMeta(encoding int, compress int) {
	defer func() {
		err := c.buf.Flush()
		if err != nil {
			panic(err)
		}
	}()

	data := make([]byte, 2)
	data[0] = utils.IntToBytes(encoding)[:1][0]
	data[1] = utils.IntToBytes(compress)[:1][0]

	c.buf.Write(data)
}

func (c TlvCodec) WriteWithLength(data []byte) {
	defer func() {
		c.buf.Flush()
		c.mu.Unlock()
	}()

	c.mu.Lock()

	l := len(data)
	lb := utils.IntToBytes(l)

	data = utils.ConcatBytes(lb, data)
	fmt.Println(" WriteWithLength len:", len(data))
	c.buf.Write(data)
}

func (c TlvCodec) Write(data interface{}) {
	c.mu.Lock()
	defer func() {
		c.buf.Flush()
		c.mu.Unlock()
	}()

	if d, ok := data.([]byte); ok {
		c.buf.Write(d)
	} else {
		panic("only accept []byte")
	}
}

func (c TlvCodec) Close() {
	(*c.conn).Close()
}

func (c TlvCodec) WriteHeader(header proto.Message) error {
	defer func() {
		c.buf.Flush()
	}()

	data, err := proto.Marshal(header)
	if err != nil {
		return err
	}

	c.buf.Write(utils.IntToBytes(len(data)))
	c.buf.Write(data)

	fmt.Println(data)

	return nil
}

func (c TlvCodec) readMeta() (*meta, error) {
	buf := make([]byte, 2)

	if n, err := (*(c.conn)).Read(buf); err != nil {
		return nil, err
	} else {
		if n != 2 {
			return nil, errors.New("not enough data to read meta")
		}
		meta := meta{
			utils.BytesToInt(buf[:1]),
			utils.BytesToInt(buf[1:]),
		}
		return &meta, nil
	}
}

func (c TlvCodec) readHeader() (*protocol.RPCHeader, error) {
	lengthBuf := make([]byte, 4)

	if n, err := (*(c.conn)).Read(lengthBuf); err != nil {
		return nil, err
	} else {
		if n != 4 {
			return nil, errors.New("not enough data to read headerLength")
		}
	}

	headerLength := utils.BytesToInt(lengthBuf[:4])

	fmt.Println("get header headerLength:", headerLength)
	data := make([]byte, headerLength)
	reader := bufio.NewReader(*c.conn)
	n, err := reader.Read(data)
	if err != nil {
		return nil, err
	}
	if n != headerLength {
		return nil, errors.New("not enough data to read header")
	}

	header := protocol.RPCHeader{}
	proto.Unmarshal(data, &header)

	return &header, nil
}

func (c TlvCodec) ReadBody() (*protocol.RPCBody, error) {
	lengthBuf := make([]byte, 4)

	fmt.Println("read before")
	if n, err := (*(c.conn)).Read(lengthBuf); err != nil {
		fmt.Println("read err:", err, " n:", n)
		return nil, err
	} else {
		fmt.Println("read body n:", n)
		if n != 4 {
			return nil, errors.New("not enough data to read headerLength")
		}
	}
	fmt.Println("read end")

	bodyLength := utils.BytesToInt(lengthBuf[:4])

	data := make([]byte, bodyLength)
	reader := bufio.NewReader(*c.conn)

	n, err := reader.Read(data)
	if err != nil {
		return nil, err
	}
	if n != bodyLength {
		return nil, errors.New("not enough data to read header")
	}

	body := protocol.RPCBody{}
	err = proto.Unmarshal(data, &body)
	if err != nil {
		return nil, err
	}

	fmt.Println("resolve tlv success")

	return &body, nil
}

func (c TlvCodec) ParseRequest(handler *serviceHandler.ServiceHandler) (*service.ServiceCall, error) {
	meta, err := c.readMeta()
	if err != nil {
		return nil, err
	}
	fmt.Println("request meta encoding:", meta.encoding)

	time.Sleep(time.Duration(1) * time.Second)

	header, err := c.readHeader()
	if err != nil {
		return nil, err
	}
	serviceMethod := header.Service
	fmt.Println("serviceMethod:", serviceMethod)

	body, err := c.ReadBody()
	if err != nil {
		fmt.Println("read body err:", err)
		return nil, err
	}

	err, method := handler.ResolveServiceMethod(serviceMethod)
	if err != nil {
		return nil, err
	}

	args := method.UnmarshalArgs(body.Args)

	call := service.NewServiceCall(method, args, method.NewReply(), int64(header.Seq))
	return &call, nil
}

func (c TlvCodec) ParseResponse() (*protocol.RPCResponseHeader, error) {
	lengthBuf := make([]byte, 4)

	if n, err := (*(c.conn)).Read(lengthBuf); err != nil {
		return nil, err
	} else {
		if n != 4 {
			return nil, errors.New("not enough data to read headerLength")
		}
	}

	headerLength := utils.BytesToInt(lengthBuf[:4])

	fmt.Println("parse headerLength:", headerLength)
	data := make([]byte, headerLength)
	reader := bufio.NewReader(*c.conn)
	n, err := reader.Read(data)
	if err != nil {
		return nil, err
	}
	if n != headerLength {
		return nil, errors.New("not enough data to read header")
	}

	header := protocol.RPCResponseHeader{}
	proto.Unmarshal(data, &header)

	return &header, nil
}
