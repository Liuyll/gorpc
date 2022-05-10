package message

type ClientHeader struct {
	Type string
	Seq int64
	Error error
	Reply interface{}
}

func NewShakeClientHeader() *ClientHeader {
	return &ClientHeader{
		Type: "acceptHandShake",
	}
}