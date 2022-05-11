package client

type Pool = map[string]*Client
type Center struct {
	serviceClientPool Pool
}

func (c Center) Call(method string, args interface{}, reply interface{}) {
	address := "127.0.0.1:8765"

	var client *Client
	if client = c.serviceClientPool[address]; client == nil {
		client = newClient(tcp, address, GobType, 1)
	}

	client.Call(method, args, reply)
}

func NewClient() *Center {
	return new(Center)
}
