package service

type ServiceCall struct {
	Method *MethodType
	Args interface{}
	Reply interface{}
	Seq int64
}

func NewServiceCall(method *MethodType, args interface{}, reply interface{}, seq int64) ServiceCall {
	return ServiceCall{
		method,
		args,
		reply,
		seq,
	}
}