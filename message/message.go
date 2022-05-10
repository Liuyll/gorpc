package message

type RPCMessage struct {
	H    *RPCHeader
	Body *RPCBody
}
