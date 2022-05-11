package client

import (
	"bufio"
	"errors"
	"fmt"
	"gorpc/message"
	"gorpc/utils"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	mu      *sync.Mutex
	maxSeq  int64
	callMap sync.Map
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

func (c *Client) handleResponse() {
	for {
		h := new(message.ClientHeader)

		if err := (*(c.codec)).ReadHeader(h); err != nil {
			if err != io.EOF {
				fmt.Println("client handleResponse readHeader err:", err)
				return
			}
		} else {
			if h.Type == "serviceResponse" {
				seq := h.Seq
				v, ok := c.callMap.Load(seq)
				if !ok {
					fmt.Println("callmap load error")
					return
				}

				if call, ok := v.(*Call); ok {
					if h.Error != "" {
						call.Error = errors.New(fmt.Sprintf("SERVER ERROR: %s", h.Error))
						call.Done()
						return
					}

					utils.SetInterfacePtr(call.Reply, h.Reply)
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

	(*(c.codec)).Write(*header, *body)
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
