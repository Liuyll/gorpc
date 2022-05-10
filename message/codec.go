package message

type Codec interface {
	Close()
	ReadHeader(interface{}) error
	ReadBody(interface{}) error
	Write(interface{}, interface{})
}
