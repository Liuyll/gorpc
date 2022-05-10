package message

type RPCHeader struct {
	ServiceMethod string
	Seq int64
	Timeout int
}
