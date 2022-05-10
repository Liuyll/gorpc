package message

type Codec interface {
	Close()
	ReadHeader(*RPCHeader) error
	ReadBody(interface{}) error
	Write(*RPCHeader, interface{})
}
