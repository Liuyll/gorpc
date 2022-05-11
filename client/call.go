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
	done chan int
	ErrorType EncodeType
}

func (c Call) WaitUntilDone() {
	<- c.done
	return
}

func (c Call) Done() {
	c.done <- 1
}
