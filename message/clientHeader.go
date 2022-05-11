package message

type ClientHeader struct {
	Type string
	Seq int64
	Error string
	Reply interface{}
}

func NewShakeClientHeader() *ClientHeader {
	return &ClientHeader{
		Type: "acceptHandShake",
	}
}