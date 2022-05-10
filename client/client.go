package client

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"gorpc/message"
	"gorpc/test"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	mu *sync.Mutex
	maxSeq int64
	callMap map[int]*Call
	codec *message.Codec
}

type CompressType int
type EncodeType int
const (
	GobType EncodeType = 1
	JsonType EncodeType = 2
)

func NewClient(encodeType EncodeType, compressType CompressType) *Client {
	gob.Register(test.Args{})

	conn, err := net.Dial("tcp", "127.0.0.1:8765")
	if err != nil {
		fmt.Println("dial err:", err)
		return nil
	}

	writer := bufio.NewWriter(conn)
	var compressTypeFlag byte = 'C'
	var encodeTypeFlag byte = 'A'
	if _, err := writer.Write([]byte{encodeTypeFlag, '_', compressTypeFlag, '_', '\n'}); err != nil {
		fmt.Println("write err:", err)
		return nil
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

	c := new(Client)
	c.codec = &codec

	c.mu = new(sync.Mutex)
	c.maxSeq = 0

	go c.handleResponse()
	return c
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
		Seq: call.Seq,
		ServiceMethod: serviceMethod,
	}

	b := message.RPCBody{
		Args: args,
		Reply: reply,
	}

	go c.innerRequest(&h, &b)

	return call
}