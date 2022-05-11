package serviceHandler

import (
	"errors"
	"fmt"
	"gorpc/service"
	"reflect"
	"strings"
)

type ServiceHandler struct {
	serviceMap map[string]*service.Service
}

func NewServiceHandler() ServiceHandler {
	return ServiceHandler{
		make(map[string]*service.Service),
	}
}

func (handler ServiceHandler) Register(maybe ...interface{}) {
	if len(maybe) == 0 {
		fmt.Printf("why send 0 argument in register?")
		return
	} else if len(maybe) == 1 {
		m0 := maybe[0]
		if reflect.TypeOf(m0).Kind() == reflect.Ptr {
			m0 = reflect.ValueOf(m0).Elem().Interface()
		}

		if s, ok := m0.(service.Service); !ok {
			fmt.Println("register one argument must be service")
			return
		} else {
			serviceName := reflect.TypeOf(s).Name()
			handler.serviceMap[serviceName] = &s
		}
	} else if len(maybe) == 2 {
		if name, ok := maybe[0].(string); !ok {
			fmt.Println("error 2")
			return
		} else {
			m1 := maybe[1]
			if reflect.TypeOf(m1).Kind() == reflect.Ptr {
				m1 = reflect.ValueOf(m1).Elem().Interface()
			}

			if s, ok := m1.(service.Service); !ok {
				fmt.Println("error 3")
				return
			} else {
				handler.serviceMap[name] = &s
			}
		}
	}
}

func (handler ServiceHandler) ResolveServiceMethod(serviceMethod string) (error, *service.MethodType) {
	m1 := strings.Split(serviceMethod, ".")
	if len(m1) != 2 {
		return errors.New("call method name is error"), nil
	}

	serviceName, serviceFunc := m1[0], m1[1]
	service := handler.serviceMap[serviceName]

	if service == nil {
		return errors.New(fmt.Sprintf("not exist this service: %s \n", serviceName)), nil
	}

	method := service.GetMethod(serviceFunc)

	if method == nil {
		return errors.New(fmt.Sprintf("not exist this method: %s \n", serviceFunc)), nil
	}

	return nil, method
}

func (handler ServiceHandler) Call(method *service.MethodType, args interface{}, reply interface{}) {
	method.Call(args, reply)
}
