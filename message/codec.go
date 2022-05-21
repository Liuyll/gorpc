package message

import (
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/serviceProto/protocol/protocol"
)

type Codec interface {
	Close()
	Write(interface{})
	ParseRequest(handler *serviceHandler.ServiceHandler) (*service.ServiceCall, error)
	WriteWithLength(data []byte)
}

type ClientCodec interface {
	Close()
	Write(interface{})
	ParseResponse() (*protocol.RPCResponseHeader, error)
	WriteWithLength(data []byte)
}
