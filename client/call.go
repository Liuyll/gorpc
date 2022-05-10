package client

type ErrorType int
const (
	TransferError ErrorType = 1
	ServiceError ErrorType = 2
)
type Call struct {
	Seq int64
	Args interface{}
	Reply interface{}
	Error error
	Done chan int
	ErrorType EncodeType
}

func (c Call) IsDone() int {
	return <- c.Done
}
