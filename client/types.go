package client

type CompressType = int
type EncodeType = int
type CommProtocol = string

const (
	tcp  CommProtocol = "tcp"
	http CommProtocol = "http"
)

const (
	GobType  EncodeType = 1
	JsonType EncodeType = 2
)
