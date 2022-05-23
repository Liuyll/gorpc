package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gorpc/message"
	"gorpc/serviceProto/protocol/protocol"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	mu      *sync.Mutex
	maxSeq  int64
	callMap sync.Map
	codec   *message.ClientCodec
	conn    io.ReadWriteCloser
}

func newClient(comm CommProtocol, address string, encodeType EncodeType, compressType CompressType) *Client {
	c := new(Client)

	connectDone := make(chan int, 1)
	c.connectToService(comm, address, encodeType, connectDone)
	<-connectDone

	c.mu = new(sync.Mutex)
	c.maxSeq = 0

	go c.handleResponse()
	return c
}

func (c *Client) connectToService(comm CommProtocol, address string, encodeType EncodeType, done chan int) {
	conn, err := net.Dial(comm, address)

	if err != nil {
		fmt.Println("dial err:", err)
		done <- 10001
	}

	c.conn = conn

	writer := bufio.NewWriter(conn)
	var compressTypeFlag byte = 'C'
	var encodeTypeFlag byte = 'A'
	if _, err := writer.Write([]byte{encodeTypeFlag, '_', compressTypeFlag, '_', '\n'}); err != nil {
		fmt.Println("write err:", err)
		done <- 10002
	}
	writer.Flush()

	var ioWriterAdapter io.ReadWriteCloser = conn
	var codec message.ClientCodec = *(message.NewTlvCodec(&ioWriterAdapter))
	if encodeType == GobType {
		//codec = message.NewGobCodec(conn)
	}

	h, err := codec.ParseResponse()
	if err != nil {
		panic("ParseResponse err:" + err.Error())
	}

	if h.Type != "acceptHandShake" {
		panic("server refuse")
	}

	fmt.Println("handle shake success")
	c.codec = &codec
	done <- 0
}

func (c *Client) handleResponse() {
	for {
		if h, err := (*(c.codec)).ParseResponse(); err != nil {
			if err != io.EOF {
				fmt.Println("client handleResponse readHeader err:", err)
				return
			}
		} else {
			if h.Type == "serviceResponse" {
				seq := h.Seq
				v, ok := c.callMap.Load(seq)
				if !ok {
					fmt.Printf("Seq: %d is not exist \n", seq)
					return
				}

				if call, ok := v.(*Call); ok {
					if h.Error != "" {
						call.Error = errors.New(fmt.Sprintf("SERVER ERROR: %s", h.Error))
						call.Done()
						return
					}

					if r, ok := call.Reply.(proto.Message); !ok {
						call.Error = errors.New("reply type is not message")
					} else {
						if err := proto.Unmarshal(h.Reply, r); err != nil {
							call.Error = err
						}
					}

					call.Done()
				} else {
					fmt.Println("call is unexpected")
					return
				}
			} else {
				fmt.Println("h.Type:", h.Type)
			}
		}
	}
}

func (c Client) innerRequest(header *message.RPCHeader, body *message.RPCBody) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//(*(c.codec)).Write(*header, *body)
}

func (c Client) innerRequestWithPb(header *protocol.RPCHeader, body *protocol.RPCBody) {
	c.mu.Lock()
	defer c.mu.Unlock()

	(*(c.codec)).WriteHeader(header)
	(*(c.codec)).WriteHeader(body)
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	curSeq := atomic.AddInt64(&c.maxSeq, 1)

	call := new(Call)
	call.Seq = curSeq
	call.Reply = reply
	call.done = make(chan int, 1)

	h := message.RPCHeader{
		Seq:           call.Seq,
		ServiceMethod: serviceMethod,
	}

	b := message.RPCBody{
		Args: args,
	}

	c.callMap.Store(curSeq, call)

	go c.innerRequest(&h, &b)
	call.WaitUntilDone()

	return call.Error
}

func (c *Client) CallWithTlv(serviceMethod string, args []byte, reply interface{}) error {
	curSeq := atomic.AddInt64(&c.maxSeq, 1)
	tlvCodec := *(c.codec)

	call := new(Call)
	call.Seq = curSeq
	call.done = make(chan int, 1)
	call.Reply = reply

	c.callMap.Store(curSeq, call)

	tlvCodec.Write([]byte{1, 1})

	h := protocol.RPCHeader{
		Encoding: 1,
		Compress: 1,
		Service:  serviceMethod,
		Seq:      int32(curSeq),
		Timeout:  0,
	}
	b := protocol.RPCBody{
		Args: args,
	}

	c.innerRequestWithPb(&h, &b)

	//tlvCodec.WriteWithLength([]byte{8, 1, 16, 1, 26, 13, 116, 101, 115, 116, 46, 65, 100, 100, 80, 114, 111, 116, 111, 32, 1})
	//tlvCodec.WriteWithLength([]byte{10, 4, 8, 1, 16, 2})

	call.WaitUntilDone()
	if call.Error != nil {
		return call.Error
	}

	return nil
}
