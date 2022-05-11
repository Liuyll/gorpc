package client

import (
	"bufio"
	"fmt"
	"gorpc/message"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	mu      *sync.Mutex
	maxSeq  int64
	callMap map[int]*Call
	codec   *message.Codec
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

	writer := bufio.NewWriter(conn)
	var compressTypeFlag byte = 'C'
	var encodeTypeFlag byte = 'A'
	if _, err := writer.Write([]byte{encodeTypeFlag, '_', compressTypeFlag, '_', '\n'}); err != nil {
		fmt.Println("write err:", err)
		done <- 10002
	}
	writer.Flush()

	var codec message.Codec
	if encodeType == GobType {
		codec = message.NewGobCodec(conn)
	}

	var h = new(message.ClientHeader)
	if err := codec.ReadHeader(h); err != nil {
		if err != io.EOF {
			panic("read conn err:" + err.Error())
		} else {
			if h.Type != "acceptHandShake" {
				panic("server refuse")
			}
		}
	}

	fmt.Println("handle shake success")
	c.codec = &codec
	done <- 0
}

func (c Client) handleResponse() {
	for {
		h := new(message.ClientHeader)

		if err := (*(c.codec)).ReadHeader(h); err != nil {
			if err != io.EOF {
				fmt.Println("client handleResponse readHeader err:", err)
			}
		} else {
			if h.Type == "serviceResponse" {
				if v, ok := h.Reply.(int); ok {
					fmt.Println("get value:", v)
				}
			} else {
				fmt.Println("h.Type:", h.Type)
			}
		}
	}
}

func (c Client) innerRequest(header *message.RPCHeader, body *message.RPCBody) {
	defer c.mu.Unlock()
	c.mu.Lock()

	(*(c.codec)).Write(*header, *body)
}

func (c Client) Call(serviceMethod string, args interface{}, reply interface{}) *Call {
	atomic.AddInt64(&c.maxSeq, 1)

	call := new(Call)
	call.Seq = c.maxSeq
	call.Reply = reply

	h := message.RPCHeader{
		Seq:           call.Seq,
		ServiceMethod: serviceMethod,
	}

	b := message.RPCBody{
		Args: args,
	}

	go c.innerRequest(&h, &b)

	return call
}
