package message

import (
	"github.com/golang/protobuf/proto"
	"gorpc/service"
	"gorpc/serviceHandler"
	"gorpc/serviceProto/protocol/protocol"
)

type Codec interface {
	Close()
	Write(interface{})
	WriteWithLength(data []byte)
	WriteHeader(header proto.Message) error
}

type ServerCodec interface {
	Codec
	ParseRequest(handler *serviceHandler.ServiceHandler) (*service.ServiceCall, error)
}

type ClientCodec interface {
	Codec
	ParseResponse() (*protocol.RPCResponseHeader, error)
}
