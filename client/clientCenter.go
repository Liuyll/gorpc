package client

import (
	"fmt"
	"sync"
)

type Pool = map[string]*Client
type Center struct {
	serviceClientPool Pool
	connectLock       *sync.Mutex
	connectBuildMap   map[string]bool
}

func (c Center) Call(method string, args interface{}, reply interface{}) error {
	address := "127.0.0.1:8765"

	client := c.getClient(address)

	return client.Call(method, args, reply)
}

func (c Center) CallWithTlv(method string, args []byte, reply interface{}) error {
	address := "127.0.0.1:8765"

	client := c.getClient(address)

	return client.CallWithTlv(method, args, reply)
}

func (c Center) getClient(address string) *Client {
	c.connectLock.Lock()
	defer c.connectLock.Unlock()

	var client *Client
	if client = c.serviceClientPool[address]; client == nil {
		fmt.Println("create new")

		client = newClient(tcp, address, GobType, 1)
		c.serviceClientPool[address] = client
	}

	return client
}

func NewClient() *Center {
	c := new(Center)
	c.serviceClientPool = make(map[string]*Client)
	c.connectBuildMap = make(map[string]bool)
	c.connectLock = new(sync.Mutex)

	return c
}
