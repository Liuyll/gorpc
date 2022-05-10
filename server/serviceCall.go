package server

import "gorpc/service"

type serviceCall struct {
	method *service.MethodType
	args interface{}
	reply interface{}
}
