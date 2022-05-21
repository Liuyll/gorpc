package message

import (
	"github.com/golang/protobuf/proto"
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/serviceProto/protocol/protocol"
)

type Codec interface {
}

type ServerCodec interface {
	Close()
	Write(interface{})
	ParseRequest(handler *serviceHandler.ServiceHandler) (*service.ServiceCall, error)
	WriteWithLength(data []byte)
	WriteHeader(header proto.Message) error
}

type ClientCodec interface {
	Close()
	Write(interface{})
	ParseResponse() (*protocol.RPCResponseHeader, error)
	WriteWithLength(data []byte)
	WriteHeader(header proto.Message) error
}
